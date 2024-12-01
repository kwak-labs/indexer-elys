package keeper

import (
	"context"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerAmmTypes "github.com/elys-network/elys/indexer/txs/amm"
	indexerTypes "github.com/elys-network/elys/indexer/types"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k msgServer) ExitPool(goCtx context.Context, msg *types.MsgExitPool) (*types.MsgExitPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	exitCoins, err := k.Keeper.ExitPool(ctx, sender, msg.PoolId, msg.ShareAmountIn, msg.MinAmountsOut, msg.TokenOutDenom, false)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	minAmountsOut := make([]indexerTypes.Token, len(msg.MinAmountsOut))
	for i, coin := range msg.MinAmountsOut {
		minAmountsOut[i] = indexerTypes.Token{
			Amount: coin.Amount.String(),
			Denom:  coin.Denom,
		}
	}

	tokenOut := make([]indexerTypes.Token, len(exitCoins))
	for i, coin := range exitCoins {
		tokenOut[i] = indexerTypes.Token{
			Amount: coin.Amount.String(),
			Denom:  coin.Denom,
		}
	}

	indexer.QueueTransaction(ctx, indexerAmmTypes.MsgExitPool{
		Sender:        msg.Sender,
		PoolID:        msg.PoolId,
		MinAmountsOut: minAmountsOut,
		ShareAmountIn: msg.ShareAmountIn.String(),
		TokenOutDenom: msg.TokenOutDenom,
		TokenOut:      tokenOut,
	}, []string{})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgExitPoolResponse{
		TokenOut: exitCoins,
	}, nil
}
