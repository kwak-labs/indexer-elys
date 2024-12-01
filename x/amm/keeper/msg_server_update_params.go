package keeper

import (
	"context"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerAmmTypes "github.com/elys-network/elys/indexer/txs/amm"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	k.Keeper.SetParams(ctx, *msg.Params)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerAmmTypes.MsgUpdateParams{
		Authority: msg.Authority,
		Params: indexerAmmTypes.Params{
			PoolCreationFee:             msg.Params.PoolCreationFee.String(),
			SlippageTrackDuration:       msg.Params.SlippageTrackDuration,
			BaseAssets:                  msg.Params.BaseAssets,
			WeightBreakingFeeExponent:   msg.Params.WeightBreakingFeeExponent.String(),
			WeightBreakingFeeMultiplier: msg.Params.WeightBreakingFeeMultiplier.String(),
			WeightBreakingFeePortion:    msg.Params.WeightBreakingFeePortion.String(),
			WeightRecoveryFeePortion:    msg.Params.WeightRecoveryFeePortion.String(),
			ThresholdWeightDifference:   msg.Params.ThresholdWeightDifference.String(),
			AllowedPoolCreators:         msg.Params.AllowedPoolCreators,
		},
	}, []string{})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdateParamsResponse{}, nil
}
