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

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k msgServer) SwapExactAmountIn(goCtx context.Context, msg *types.MsgSwapExactAmountIn) (*types.MsgSwapExactAmountInResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Swap event is handled elsewhere
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	response, err := k.Keeper.SwapExactAmountIn(ctx, msg)
	if err != nil {
		return nil, err
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	routes := make([]indexerAmmTypes.SwapAmountInRoute, len(msg.Routes))
	for i, route := range msg.Routes {
		routes[i] = indexerAmmTypes.SwapAmountInRoute{
			PoolID:        route.PoolId,
			TokenOutDenom: route.TokenOutDenom,
		}
	}

	indexer.QueueTransaction(ctx, indexerAmmTypes.MsgSwapExactAmountIn{
		Sender: msg.Sender,
		Routes: routes,
		TokenIn: indexerTypes.Token{
			Amount: msg.TokenIn.Amount.String(),
			Denom:  msg.TokenIn.Denom,
		},
		TokenOutMinAmount: msg.TokenOutMinAmount.String(),
		Recipient:         msg.Recipient,
		SwapFee:           response.SwapFee.String(),
		Discount:          response.Discount.String(),
		TokenOut: indexerTypes.Token{
			Amount: response.TokenOutAmount.String(),
			Denom:  msg.Routes[len(msg.Routes)-1].TokenOutDenom,
		},
	}, []string{msg.Recipient})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return response, nil
}

func (k Keeper) SwapExactAmountIn(ctx sdk.Context, msg *types.MsgSwapExactAmountIn) (*types.MsgSwapExactAmountInResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		recipient = sender
	}
	// Try executing the tx on cached context environment, to filter invalid transactions out
	cacheCtx, _ := ctx.CacheContext()
	tokenOutAmount, swapFee, discount, err := k.RouteExactAmountIn(cacheCtx, sender, recipient, msg.Routes, msg.TokenIn, sdkmath.Int(msg.TokenOutMinAmount))
	if err != nil {
		return nil, err
	}

	lastSwapIndex := k.GetLastSwapRequestIndex(ctx)
	k.SetSwapExactAmountInRequests(ctx, msg, lastSwapIndex+1)
	k.SetLastSwapRequestIndex(ctx, lastSwapIndex+1)

	return &types.MsgSwapExactAmountInResponse{
		TokenOutAmount: tokenOutAmount,
		SwapFee:        swapFee,
		Discount:       discount,
		Recipient:      recipient.String(),
	}, nil
}
