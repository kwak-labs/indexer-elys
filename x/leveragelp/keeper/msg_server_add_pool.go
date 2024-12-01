package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerLeveragelpTypes "github.com/elys-network/elys/indexer/txs/leveragelp"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) AddPool(goCtx context.Context, msg *types.MsgAddPool) (*types.MsgAddPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ammPool, ammFound := k.amm.GetPool(ctx, msg.Pool.AmmPoolId)

	if !ammFound {
		return nil, fmt.Errorf("amm pool not found")
	}

	if !ammPool.PoolParams.UseOracle {
		return nil, fmt.Errorf("amm pool does not use oracle")
	}

	_, found := k.GetPool(ctx, msg.Pool.AmmPoolId)
	if found {
		return nil, fmt.Errorf("pool already exists")
	}

	maxLeverageAllowed := k.GetMaxLeverageParam(ctx)
	leverage := sdkmath.LegacyMinDec(msg.Pool.LeverageMax, maxLeverageAllowed)

	newPool := types.NewPool(ammPool.PoolId, leverage)
	k.SetPool(ctx, newPool)

	if k.hooks != nil {
		err := k.hooks.AfterEnablingPool(ctx, ammPool)
		if err != nil {
			return nil, err
		}
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerLeveragelpTypes.MsgAddPool{
		Authority: msg.Authority,
		Pool: indexerLeveragelpTypes.AddPool{
			AmmPoolID:   msg.Pool.AmmPoolId,
			LeverageMax: msg.Pool.LeverageMax.String(),
			Leverage:    leverage.String(),
		},
	}, []string{})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgAddPoolResponse{}, nil
}
