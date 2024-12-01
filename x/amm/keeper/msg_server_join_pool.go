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

func (k msgServer) JoinPool(goCtx context.Context, msg *types.MsgJoinPool) (*types.MsgJoinPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	neededLp, sharesOut, err := k.Keeper.JoinPoolNoSwap(ctx, sender, msg.PoolId, msg.ShareAmountOut, msg.MaxAmountsIn)
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
	maxAmountsIn := make([]indexerTypes.Token, len(msg.MaxAmountsIn))
	for i, coin := range msg.MaxAmountsIn {
		maxAmountsIn[i] = indexerTypes.Token{
			Amount: coin.Amount.String(),
			Denom:  coin.Denom,
		}
	}

	tokenIn := make([]indexerTypes.Token, len(neededLp))
	for i, coin := range neededLp {
		tokenIn[i] = indexerTypes.Token{
			Amount: coin.Amount.String(),
			Denom:  coin.Denom,
		}
	}

	indexer.QueueTransaction(ctx, indexerAmmTypes.MsgJoinPool{
		Sender:         msg.Sender,
		PoolID:         msg.PoolId,
		MaxAmountsIn:   maxAmountsIn,
		ShareAmountOut: sharesOut.String(),
		TokenIn:        tokenIn,
	}, []string{})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgJoinPoolResponse{
		ShareAmountOut: sharesOut,
		TokenIn:        neededLp,
	}, nil
}
