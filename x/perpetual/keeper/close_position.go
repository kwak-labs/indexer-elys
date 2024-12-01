package keeper

import (
	"fmt"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/

	indexer "github.com/elys-network/elys/indexer"
	indexerPerpetualTypes "github.com/elys-network/elys/indexer/txs/perpetual"
	indexerTypes "github.com/elys-network/elys/indexer/types"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) ClosePosition(ctx sdk.Context, msg *types.MsgClose, baseCurrency string) (*types.MTP, math.Int, math.LegacyDec, error) {
	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	var initialCollateral math.LegacyDec
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	// Retrieve MTP
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	mtp, err := k.GetMTP(ctx, creator, msg.Id)
	if err != nil {
		return nil, math.ZeroInt(), math.LegacyZeroDec(), err
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	initialCollateral = math.LegacyNewDecFromInt(mtp.Collateral)
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	// Retrieve AmmPool
	ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId)
	if err != nil {
		return nil, math.ZeroInt(), math.LegacyZeroDec(), err
	}

	// This needs to be updated here to check user doesn't send more than required amount
	k.UpdateMTPBorrowInterestUnpaidLiability(ctx, &mtp)
	// Retrieve Pool
	pool, found := k.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return nil, math.ZeroInt(), math.LegacyZeroDec(), errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", mtp.AmmPoolId))
	}

	// Handle Borrow Interest if within epoch position SettleMTPBorrowInterestUnpaidLiability settles interest using mtp.Custody, mtp.Custody gets reduced
	if _, err = k.SettleMTPBorrowInterestUnpaidLiability(ctx, &mtp, &pool, ammPool); err != nil {
		return nil, math.ZeroInt(), math.LegacyZeroDec(), err
	}

	err = k.SettleFunding(ctx, &mtp, &pool, ammPool)
	if err != nil {
		return nil, math.ZeroInt(), math.LegacyZeroDec(), errorsmod.Wrapf(err, "error handling funding fee")
	}

	// Should be declared after SettleMTPBorrowInterestUnpaidLiability and settling funding
	closingRatio := msg.Amount.ToLegacyDec().Quo(mtp.Custody.ToLegacyDec())
	if mtp.Position == types.Position_SHORT {
		closingRatio = msg.Amount.ToLegacyDec().Quo(mtp.Liabilities.ToLegacyDec())
	}
	if closingRatio.GT(math.LegacyOneDec()) {
		closingRatio = math.LegacyOneDec()
	}

	// Estimate swap and repay
	repayAmt, err := k.EstimateAndRepay(ctx, &mtp, &pool, &ammPool, baseCurrency, closingRatio)
	if err != nil {
		return nil, math.ZeroInt(), math.LegacyZeroDec(), err
	}

	// EpochHooks after perpetual position closed
	if k.hooks != nil {
		params := k.GetParams(ctx)
		err = k.hooks.AfterPerpetualPositionClosed(ctx, ammPool, pool, creator, params.EnableTakeProfitCustodyLiabilities)
		if err != nil {
			return nil, math.Int{}, math.LegacyDec{}, err
		}
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	finalValue := math.LegacyNewDecFromInt(repayAmt).Sub(mtp.Liabilities.ToLegacyDec())
	profitLoss, profitLossPerc := calculateProfitLoss(initialCollateral, finalValue)

	indexer.QueueTransaction(ctx, indexerPerpetualTypes.MsgClose{
		Creator:  msg.Creator,
		Id:       msg.Id,
		Amount:   msg.Amount.String(),
		Position: mtp.Position.String(),
		Collateral: indexerTypes.Token{
			Amount: mtp.Collateral.String(),
			Denom:  mtp.CollateralAsset,
		},
		Custody: indexerTypes.Token{
			Amount: mtp.Custody.String(),
			Denom:  mtp.CustodyAsset,
		},
		Liabilities: indexerTypes.Token{
			Amount: mtp.Liabilities.String(),
			Denom:  mtp.LiabilitiesAsset,
		},
		RepayAmount:      repayAmt.String(),
		InitialValue:     initialCollateral.String(),
		FinalValue:       finalValue.String(),
		ProfitLoss:       profitLoss.String(),
		ProfitLossPerc:   profitLossPerc.String(),
		CollateralAsset:  mtp.CollateralAsset,
		TradingAsset:     mtp.TradingAsset,
		LiabilitiesAsset: mtp.LiabilitiesAsset,
		MtpHealth:        mtp.MtpHealth.String(),
		OpenPrice:        mtp.OpenPrice.String(),
	}, []string{msg.Creator})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &mtp, repayAmt, closingRatio, nil
}

/* *************************************************************************** */
/* Start of kwak-indexer node implementation*/
func calculateProfitLoss(
	initialValue math.LegacyDec,
	finalValue math.LegacyDec,
) (profitLoss math.LegacyDec, profitLossPerc math.LegacyDec) {
	profitLoss = finalValue.Sub(initialValue)
	if !initialValue.IsZero() {
		profitLossPerc = profitLoss.Quo(initialValue).Mul(math.LegacyNewDec(100))
	}
	return profitLoss, profitLossPerc
}

/* End of kwak-indexer node implementation*/
/* *************************************************************************** */
