package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerAmmTypes "github.com/elys-network/elys/indexer/txs/amm"
	indexerTypes "github.com/elys-network/elys/indexer/types"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k msgServer) SwapByDenom(goCtx context.Context, msg *types.MsgSwapByDenom) (*types.MsgSwapByDenomResponse, error) {
	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	// ! this entire function was changed a bit
	// ! Originally it was simply return k.Keeper.SwapByDenom(ctx, msg)
	// ! But SwapByDenom is called a lot in mutiple instances not from TXs, so to maintain simplicity of being able to indentify exact swap outputs from a tx I did it this way
	ctx := sdk.UnwrapSDKContext(goCtx)

	response, err := k.Keeper.SwapByDenom(ctx, msg)
	if err != nil {
		return nil, err
	}

	var indexerInRoute []indexerAmmTypes.SwapAmountInRoute
	if response.InRoute != nil {
		indexerInRoute = make([]indexerAmmTypes.SwapAmountInRoute, len(response.InRoute))
		for i, route := range response.InRoute {
			indexerInRoute[i] = indexerAmmTypes.SwapAmountInRoute{
				PoolID:        route.PoolId,
				TokenOutDenom: route.TokenOutDenom,
			}
		}
	}

	// Convert outRoute to indexer format
	var indexerOutRoute []indexerAmmTypes.SwapAmountOutRoute
	if response.OutRoute != nil {
		indexerOutRoute = make([]indexerAmmTypes.SwapAmountOutRoute, len(response.OutRoute))
		for i, route := range response.OutRoute {
			indexerOutRoute[i] = indexerAmmTypes.SwapAmountOutRoute{
				PoolID:       route.PoolId,
				TokenInDenom: route.TokenInDenom,
			}
		}
	}

	indexer.QueueTransaction(ctx, indexerAmmTypes.MsgSwapByDenom{
		Sender: msg.Sender,
		Amount: indexerTypes.Token{
			Amount: msg.Amount.Amount.String(),
			Denom:  msg.Amount.Denom,
		},
		MinAmount: indexerTypes.Token{
			Amount: msg.MinAmount.Amount.String(),
			Denom:  msg.MinAmount.Denom,
		},
		MaxAmount: indexerTypes.Token{
			Amount: msg.MaxAmount.Amount.String(),
			Denom:  msg.MaxAmount.Denom,
		},
		DenomIn:   msg.DenomIn,
		DenomOut:  msg.DenomOut,
		Recipient: msg.Recipient,
		InRoute:   indexerInRoute,
		OutRoute:  indexerOutRoute,
		SpotPrice: response.SpotPrice.String(),
		SwapFee:   response.SwapFee.String(),
		Discount:  response.Discount.String(),
		TokenOut: indexerTypes.Token{
			Amount: response.Amount.Amount.String(),
			Denom:  response.Amount.Denom,
		},
	}, []string{msg.Sender, msg.Recipient})
	return response, nil
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

}

func (k Keeper) SwapByDenom(ctx sdk.Context, msg *types.MsgSwapByDenom) (*types.MsgSwapByDenomResponse, error) {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// retrieve base currency denom
	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	inRoute, outRoute, _, spotPrice, _, _, _, _, _, _, err := k.CalcSwapEstimationByDenom(ctx, msg.Amount, msg.DenomIn, msg.DenomOut, baseCurrency, msg.Sender, sdkmath.LegacyZeroDec(), 0)
	if err != nil {
		return nil, err
	}

	// swap to token out with exact amount in using in route
	if inRoute != nil {
		// check min amount denom is equals to denom out
		if msg.MinAmount.Denom != msg.DenomOut {
			return nil, errorsmod.Wrapf(types.ErrInvalidDenom, "min amount denom %s is not equals to denom out %s", msg.MinAmount.Denom, msg.DenomOut)
		}

		// convert route []*types.SwapAmountInRoute to []types.SwapAmountInRoute
		route := make([]types.SwapAmountInRoute, len(inRoute))
		for i, r := range inRoute {
			route[i] = *r
		}

		res, err := k.SwapExactAmountIn(
			ctx,
			&types.MsgSwapExactAmountIn{
				Sender:            msg.Sender,
				Recipient:         msg.Recipient,
				Routes:            route,
				TokenIn:           msg.Amount,
				TokenOutMinAmount: msg.MinAmount.Amount,
			},
		)
		if err != nil {
			return nil, err
		}

		return &types.MsgSwapByDenomResponse{
			Amount:    sdk.NewCoin(msg.DenomOut, res.TokenOutAmount),
			InRoute:   inRoute,
			OutRoute:  nil,
			SpotPrice: spotPrice,
			SwapFee:   res.SwapFee,
			Discount:  res.Discount,
			Recipient: res.Recipient,
		}, nil
	}

	// swap to token in with exact amount out using out route
	if outRoute != nil {
		// check max amount denom is equals to denom out
		if msg.MaxAmount.Denom != msg.DenomOut {
			return nil, errorsmod.Wrapf(types.ErrInvalidDenom, "max amount denom %s is not equals to denom out %s", msg.MaxAmount.Denom, msg.DenomOut)
		}

		// convert route []*types.SwapAmountOutRoute to []types.SwapAmountOutRoute
		route := make([]types.SwapAmountOutRoute, len(outRoute))
		for i, r := range outRoute {
			route[i] = *r
		}

		res, err := k.SwapExactAmountOut(
			ctx,
			&types.MsgSwapExactAmountOut{
				Sender:           msg.Sender,
				Routes:           route,
				TokenInMaxAmount: msg.MaxAmount.Amount,
				TokenOut:         msg.Amount,
			},
		)
		if err != nil {
			return nil, err
		}

		return &types.MsgSwapByDenomResponse{
			Amount:    sdk.NewCoin(msg.DenomOut, res.TokenInAmount),
			InRoute:   nil,
			OutRoute:  outRoute,
			SpotPrice: spotPrice,
			SwapFee:   res.SwapFee,
			Discount:  res.Discount,
			Recipient: res.Recipient,
		}, nil
	}

	// otherwise throw an error
	return nil, errorsmod.Wrapf(types.ErrInvalidSwapMsgType, "neither inRoute nor outRoute are available")
}
