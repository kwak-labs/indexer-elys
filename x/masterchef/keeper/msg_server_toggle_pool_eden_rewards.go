package keeper

import (
	"context"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerMasterchefTypes "github.com/elys-network/elys/indexer/txs/masterchef"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (k msgServer) TogglePoolEdenRewards(goCtx context.Context, msg *types.MsgTogglePoolEdenRewards) (*types.MsgTogglePoolEdenRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	pool, found := k.GetPoolInfo(ctx, msg.PoolId)
	if !found {
		return &types.MsgTogglePoolEdenRewardsResponse{}, types.ErrPoolNotFound
	}

	pool.EnableEdenRewards = msg.Enable
	k.SetPoolInfo(ctx, pool)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerMasterchefTypes.MsgTogglePoolEdenRewards{
		Authority: msg.Authority,
		PoolID:    msg.PoolId,
		Enable:    msg.Enable,
	}, []string{msg.Authority})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgTogglePoolEdenRewardsResponse{}, nil
}
