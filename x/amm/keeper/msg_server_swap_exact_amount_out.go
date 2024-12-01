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

func (k msgServer) SwapExactAmountOut(goCtx context.Context, msg *types.MsgSwapExactAmountOut) (*types.MsgSwapExactAmountOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Swap event is handled elsewhere
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	// ! this entire function was changed a bit
	// ! Originally it was simply return k.Keeper.SwapExactAmountOut(ctx, msg)
	// ! But to avoid the keeper being called elsewhere and messing witth things, had to edit this to get the ouput of the swap in context of the tx
	response, err := k.Keeper.SwapExactAmountOut(ctx, msg)
	if err != nil {
		return nil, err
	}

	routes := make([]indexerAmmTypes.SwapAmountOutRoute, len(msg.Routes))
	for i, route := range msg.Routes {
		routes[i] = indexerAmmTypes.SwapAmountOutRoute{
			PoolID:       route.PoolId,
			TokenInDenom: route.TokenInDenom,
		}
	}

	indexer.QueueTransaction(ctx, indexerAmmTypes.MsgSwapExactAmountOut{
		Sender: msg.Sender,
		Routes: routes,
		TokenOut: indexerTypes.Token{
			Amount: msg.TokenOut.Amount.String(),
			Denom:  msg.TokenOut.Denom,
		},
		TokenInMaxAmount: msg.TokenInMaxAmount.String(),
		Recipient:        msg.Recipient,
		TokenInAmount: indexerTypes.Token{
			Amount: response.TokenInAmount.String(),
			Denom:  msg.Routes[0].TokenInDenom, // First route's token denom is the input token
		},
		SwapFee: indexerTypes.Token{
			Amount: response.SwapFee.String(),
			Denom:  msg.TokenOut.Denom, // Fee is in output token denomination
		},
		Discount: indexerTypes.Token{
			Amount: response.Discount.String(),
			Denom:  msg.TokenOut.Denom, // Discount is in output token denomination
		},
	}, []string{msg.Recipient})

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return response, nil
}

func (k Keeper) SwapExactAmountOut(ctx sdk.Context, msg *types.MsgSwapExactAmountOut) (*types.MsgSwapExactAmountOutResponse, error) {
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
	tokenInAmount, swapFee, discount, err := k.RouteExactAmountOut(cacheCtx, sender, recipient, msg.Routes, msg.TokenInMaxAmount, msg.TokenOut)
	if err != nil {
		return nil, err
	}

	lastSwapIndex := k.GetLastSwapRequestIndex(ctx)
	k.SetSwapExactAmountOutRequests(ctx, msg, lastSwapIndex+1)
	k.SetLastSwapRequestIndex(ctx, lastSwapIndex+1)

	return &types.MsgSwapExactAmountOutResponse{
		TokenInAmount: tokenInAmount,
		SwapFee:       swapFee,
		Discount:      discount,
		Recipient:     recipient.String(),
	}, nil
}
