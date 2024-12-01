package keeper

import (
	"context"
	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerParamTypes "github.com/elys-network/elys/indexer/txs/parameter"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/parameter/types"
)

func (k msgServer) UpdateRewardsDataLifetime(goCtx context.Context, msg *types.MsgUpdateRewardsDataLifetime) (*types.MsgUpdateRewardsDataLifetimeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	params := k.GetParams(ctx)
	params.RewardsDataLifetime = msg.RewardsDataLifetime
	k.SetParams(ctx, params)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerParamTypes.MsgUpdateRewardsDataLifetime{
		Creator:             msg.Creator,
		RewardsDataLifetime: msg.RewardsDataLifetime,
	}, []string{msg.Creator})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdateRewardsDataLifetimeResponse{}, nil
}
