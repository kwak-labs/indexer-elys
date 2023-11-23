package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/amm/types"
)

// CalcOutRouteByDenom calculates the out route by denom
func (k Keeper) CalcOutRouteByDenom(ctx sdk.Context, denomOut string, denomIn string, baseCurrency string) ([]*types.SwapAmountOutRoute, error) {
	var route []*types.SwapAmountOutRoute

	// If the denoms are the same, throw an error
	if denomIn == denomOut {
		return nil, sdkerrors.Wrap(types.ErrSameDenom, "denom in and denom out are the same")
	}

	// Check for a direct pool between the denoms
	if poolId, found := k.GetPoolIdWithAllDenoms(ctx, []string{denomOut, denomIn}); found {
		// If the pool exists, return the route
		route = append(route, &types.SwapAmountOutRoute{
			PoolId:       poolId,
			TokenInDenom: denomIn,
		})
		return route, nil
	}

	// Find pool for initial denom to base currency
	poolId, found := k.GetPoolIdWithAllDenoms(ctx, []string{denomOut, baseCurrency})
	if !found {
		return nil, fmt.Errorf("no available pool for %s to base currency", denomOut)
	}
	// If the pool exists, append the route
	route = append(route, &types.SwapAmountOutRoute{
		PoolId:       poolId,
		TokenInDenom: baseCurrency,
	})

	// Find pool for base currency to target denom
	poolId, found = k.GetPoolIdWithAllDenoms(ctx, []string{baseCurrency, denomIn})
	if !found {
		return nil, fmt.Errorf("no available pool for base currency to %s", denomIn)
	}
	// If the pool exists, append the route
	route = append(route, &types.SwapAmountOutRoute{
		PoolId:       poolId,
		TokenInDenom: denomIn,
	})

	return route, nil
}