package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/parameter/types"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerParamTypes "github.com/elys-network/elys/indexer/txs/parameter"
	/* End of kwak-indexer node implementation*/ /* *************************************************************************** */)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) UpdateMinCommission(goCtx context.Context, msg *types.MsgUpdateMinCommission) (*types.MsgUpdateMinCommissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	params := k.GetParams(ctx)
	params.MinCommissionRate = msg.MinCommission
	k.SetParams(ctx, params)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerParamTypes.MsgUpdateMinCommission{
		Creator:       msg.Creator,
		MinCommission: msg.MinCommission.String(),
	}, []string{msg.Creator})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdateMinCommissionResponse{}, nil
}

func (k msgServer) UpdateMaxVotingPower(goCtx context.Context, msg *types.MsgUpdateMaxVotingPower) (*types.MsgUpdateMaxVotingPowerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	params := k.GetParams(ctx)
	params.MaxVotingPower = msg.MaxVotingPower
	k.SetParams(ctx, params)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerParamTypes.MsgUpdateMaxVotingPower{
		Creator:        msg.Creator,
		MaxVotingPower: msg.MaxVotingPower.String(),
	}, []string{msg.Creator})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdateMaxVotingPowerResponse{}, nil
}

func (k msgServer) UpdateMinSelfDelegation(goCtx context.Context, msg *types.MsgUpdateMinSelfDelegation) (*types.MsgUpdateMinSelfDelegationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	params := k.GetParams(ctx)
	params.MinSelfDelegation = msg.MinSelfDelegation
	k.SetParams(ctx, params)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerParamTypes.MsgUpdateMinSelfDelegation{
		Creator:           msg.Creator,
		MinSelfDelegation: msg.MinSelfDelegation.String(),
	}, []string{msg.Creator})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdateMinSelfDelegationResponse{}, nil
}

func (k msgServer) UpdateTotalBlocksPerYear(goCtx context.Context, msg *types.MsgUpdateTotalBlocksPerYear) (*types.MsgUpdateTotalBlocksPerYearResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	params := k.GetParams(ctx)
	params.TotalBlocksPerYear = msg.TotalBlocksPerYear
	k.SetParams(ctx, params)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerParamTypes.MsgUpdateTotalBlocksPerYear{
		Creator:            msg.Creator,
		TotalBlocksPerYear: msg.TotalBlocksPerYear,
	}, []string{msg.Creator})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdateTotalBlocksPerYearResponse{}, nil
}
