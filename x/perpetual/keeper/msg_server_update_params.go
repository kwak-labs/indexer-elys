package keeper

import (
	"context"
	"fmt"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerPerpetualTypes "github.com/elys-network/elys/indexer/txs/perpetual"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// Update params through gov proposal
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// store params
	if err := k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	if k.hooks != nil {
		pools := k.GetAllPools(ctx)
		for _, pool := range pools {
			ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId)
			if err != nil {
				return nil, fmt.Errorf("amm pool %d not found", pool.AmmPoolId)
			}

			err = k.hooks.AfterParamsChange(ctx, ammPool, pool, msg.Params.EnableTakeProfitCustodyLiabilities)
			if err != nil {
				return nil, err
			}
		}
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	// Queue the params update transaction
	indexer.QueueTransaction(ctx, indexerPerpetualTypes.MsgUpdateParams{
		Authority: msg.Authority,
		Params: indexerPerpetualTypes.Params{
			LeverageMax:                                    msg.Params.LeverageMax.String(),
			BorrowInterestRateMax:                          msg.Params.BorrowInterestRateMax.String(),
			BorrowInterestRateMin:                          msg.Params.BorrowInterestRateMin.String(),
			BorrowInterestRateIncrease:                     msg.Params.BorrowInterestRateIncrease.String(),
			BorrowInterestRateDecrease:                     msg.Params.BorrowInterestRateDecrease.String(),
			HealthGainFactor:                               msg.Params.HealthGainFactor.String(),
			MaxOpenPositions:                               msg.Params.MaxOpenPositions,
			PoolOpenThreshold:                              msg.Params.PoolOpenThreshold.String(),
			ForceCloseFundPercentage:                       msg.Params.ForceCloseFundPercentage.String(),
			ForceCloseFundAddress:                          msg.Params.ForceCloseFundAddress,
			IncrementalBorrowInterestPaymentFundPercentage: msg.Params.IncrementalBorrowInterestPaymentFundPercentage.String(),
			IncrementalBorrowInterestPaymentFundAddress:    msg.Params.IncrementalBorrowInterestPaymentFundAddress,
			SafetyFactor:                                   msg.Params.SafetyFactor.String(),
			IncrementalBorrowInterestPaymentEnabled:        msg.Params.IncrementalBorrowInterestPaymentEnabled,
			WhitelistingEnabled:                            msg.Params.WhitelistingEnabled,
			PerpetualSwapFee:                               msg.Params.PerpetualSwapFee.String(),
			MaxLimitOrder:                                  msg.Params.MaxLimitOrder,
			FixedFundingRate:                               msg.Params.FixedFundingRate.String(),
			MinimumLongTakeProfitPriceRatio:                msg.Params.MinimumLongTakeProfitPriceRatio.String(),
			MaximumLongTakeProfitPriceRatio:                msg.Params.MaximumLongTakeProfitPriceRatio.String(),
			MaximumShortTakeProfitPriceRatio:               msg.Params.MaximumShortTakeProfitPriceRatio.String(),
			EnableTakeProfitCustodyLiabilities:             msg.Params.EnableTakeProfitCustodyLiabilities,
			WeightBreakingFeeFactor:                        msg.Params.WeightBreakingFeeFactor.String(),
		},
	}, []string{msg.Authority})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdateParamsResponse{}, nil
}
