package keeper

import (
	"context"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerStableStakeTypes "github.com/elys-network/elys/indexer/txs/stablestake"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

// Update params through gov proposal
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// overwrite total value in params
	params := k.GetParams(ctx)
	msg.Params.TotalValue = params.TotalValue

	// store params
	k.SetParams(ctx, *msg.Params)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerStableStakeTypes.MsgUpdateParams{
		Authority: msg.Authority,
		Params: indexerStableStakeTypes.Params{
			DepositDenom:         msg.Params.DepositDenom,
			RedemptionRate:       msg.Params.RedemptionRate.String(),
			EpochLength:          msg.Params.EpochLength,
			InterestRate:         msg.Params.InterestRate.String(),
			InterestRateMax:      msg.Params.InterestRateMax.String(),
			InterestRateMin:      msg.Params.InterestRateMin.String(),
			InterestRateIncrease: msg.Params.InterestRateIncrease.String(),
			InterestRateDecrease: msg.Params.InterestRateDecrease.String(),
			HealthGainFactor:     msg.Params.HealthGainFactor.String(),
			TotalValue:           msg.Params.TotalValue.String(),
			MaxLeverageRatio:     msg.Params.MaxLeverageRatio.String(),
		},
	}, []string{msg.Authority})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdateParamsResponse{}, nil
}
