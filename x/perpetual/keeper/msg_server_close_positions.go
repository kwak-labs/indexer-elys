package keeper

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	"cosmossdk.io/math"
	indexer "github.com/elys-network/elys/indexer"
	indexerPerpetualTypes "github.com/elys-network/elys/indexer/txs/perpetual"
	indexerTypes "github.com/elys-network/elys/indexer/types"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k msgServer) ClosePositions(goCtx context.Context, msg *types.MsgClosePositions) (*types.MsgClosePositionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	// Queue the initial close positions transaction
	liquidateRequests := make([]indexerPerpetualTypes.PositionRequest, len(msg.Liquidate))
	for i, req := range msg.Liquidate {
		liquidateRequests[i] = indexerPerpetualTypes.PositionRequest{
			Address: req.Address,
			ID:      req.Id,
		}
	}

	stopLossRequests := make([]indexerPerpetualTypes.PositionRequest, len(msg.StopLoss))
	for i, req := range msg.StopLoss {
		stopLossRequests[i] = indexerPerpetualTypes.PositionRequest{
			Address: req.Address,
			ID:      req.Id,
		}
	}

	takeProfitRequests := make([]indexerPerpetualTypes.PositionRequest, len(msg.TakeProfit))
	for i, req := range msg.TakeProfit {
		takeProfitRequests[i] = indexerPerpetualTypes.PositionRequest{
			Address: req.Address,
			ID:      req.Id,
		}
	}

	indexer.QueueTransaction(ctx, indexerPerpetualTypes.MsgClosePositions{
		Creator:    msg.Creator,
		Liquidate:  liquidateRequests,
		StopLoss:   stopLossRequests,
		TakeProfit: takeProfitRequests,
	}, []string{msg.Creator})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	baseCurrency, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, nil
	}

	// Handle liquidations
	liqLog := []string{}
	for _, val := range msg.Liquidate {
		owner := sdk.MustAccAddressFromBech32(val.Address)
		position, err := k.GetMTP(ctx, owner, val.Id)
		if err != nil {
			continue
		}

		pool, poolFound := k.GetPool(ctx, position.AmmPoolId)
		if !poolFound {
			continue
		}
		ammPool, poolErr := k.GetAmmPool(ctx, position.AmmPoolId)
		if poolErr != nil {
			continue
		}

		err = k.CheckAndLiquidateUnhealthyPosition(ctx, &position, pool, ammPool, baseCurrency.Denom)
		if err != nil {
			// Add log about error or not liquidated
			liqLog = append(liqLog, fmt.Sprintf("Position: Address:%s Id:%d cannot be liquidated due to err: %s", position.Address, position.Id, err.Error()))
		} else {
			/* *************************************************************************** */
			/* Start of kwak-indexer node implementation*/
			initialValue := math.LegacyNewDecFromInt(position.Collateral)
			// For perpetual positions, we calculate final value based on custody and liabilities
			finalValue := math.LegacyNewDecFromInt(position.Custody.Sub(position.Liabilities))
			profitLoss := finalValue.Sub(initialValue)
			var profitLossPerc math.LegacyDec
			if !initialValue.IsZero() {
				profitLossPerc = profitLoss.Quo(initialValue).Mul(math.LegacyNewDec(100))
			}

			indexer.QueueEvent(ctx, "/elys-event/perpetual/liquidation", indexerPerpetualTypes.ClosePositionEvent{
				Address: position.Address,
				ID:      position.Id,
				Collateral: indexerTypes.Token{
					Amount: position.Collateral.String(),
					Denom:  position.CollateralAsset,
				},
				Custody: indexerTypes.Token{
					Amount: position.Custody.String(),
					Denom:  position.CustodyAsset,
				},
				Liabilities: indexerTypes.Token{
					Amount: position.Liabilities.String(),
					Denom:  position.LiabilitiesAsset,
				},
				Health:         position.MtpHealth.String(),
				InitialValue:   initialValue.String(),
				FinalValue:     finalValue.String(),
				ProfitLoss:     profitLoss.String(),
				ProfitLossPerc: profitLossPerc.String(),
				OpenPrice:      position.OpenPrice.String(),
				Position:       position.Position.String(),
			}, []string{position.Address})

			/* End of kwak-indexer node implementation*/
			/* *************************************************************************** */

			ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClosePositions,
				sdk.NewAttribute("address", position.Address),
				sdk.NewAttribute("id", strconv.FormatUint(position.Id, 10)),
			))
		}
	}

	//Handle StopLoss
	closeLog := []string{}
	for _, val := range msg.StopLoss {
		owner := sdk.MustAccAddressFromBech32(val.Address)
		position, err := k.GetMTP(ctx, owner, val.Id)
		if err != nil {
			continue
		}

		pool, poolFound := k.GetPool(ctx, position.AmmPoolId)
		if !poolFound {
			continue
		}

		err = k.CheckAndCloseAtStopLoss(ctx, &position, pool, baseCurrency.Denom)
		if err != nil {
			// Add log about error or not closed
			closeLog = append(closeLog, fmt.Sprintf("Position: Address:%s Id:%d cannot be liquidated due to err: %s", position.Address, position.Id, err.Error()))
		} else {
			/* *************************************************************************** */
			/* Start of kwak-indexer node implementation*/
			initialValue := math.LegacyNewDecFromInt(position.Collateral)
			// For perpetual positions, we calculate final value based on custody and liabilities
			finalValue := math.LegacyNewDecFromInt(position.Custody.Sub(position.Liabilities))
			profitLoss := finalValue.Sub(initialValue)
			var profitLossPerc math.LegacyDec
			if !initialValue.IsZero() {
				profitLossPerc = profitLoss.Quo(initialValue).Mul(math.LegacyNewDec(100))
			}

			indexer.QueueEvent(ctx, "/elys-event/perpetual/stop-loss", indexerPerpetualTypes.StopLossEvent{
				Address:  position.Address,
				ID:       position.Id,
				Position: position.Position.String(),
				Collateral: indexerTypes.Token{
					Amount: position.Collateral.String(),
					Denom:  position.CollateralAsset,
				},
				Custody: indexerTypes.Token{
					Amount: position.Custody.String(),
					Denom:  position.CustodyAsset,
				},
				Liabilities: indexerTypes.Token{
					Amount: position.Liabilities.String(),
					Denom:  position.LiabilitiesAsset,
				},
				StopLossPrice:  position.StopLossPrice.String(),
				OpenPrice:      position.OpenPrice.String(),
				Health:         position.MtpHealth.String(),
				ProfitLoss:     profitLoss.String(),
				ProfitLossPerc: profitLossPerc.String(),
			}, []string{position.Address})
			/* End of kwak-indexer node implementation*/
			/* *************************************************************************** */

			ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClosePositions,
				sdk.NewAttribute("address", position.Address),
				sdk.NewAttribute("id", strconv.FormatUint(position.Id, 10)),
			))
		}
	}

	//Handle take profit
	takeProfitLog := []string{}
	for _, val := range msg.TakeProfit {
		owner := sdk.MustAccAddressFromBech32(val.Address)
		position, err := k.GetMTP(ctx, owner, val.Id)
		if err != nil {
			continue
		}

		pool, poolFound := k.GetPool(ctx, position.AmmPoolId)
		if !poolFound {
			continue
		}

		err = k.CheckAndCloseAtTakeProfit(ctx, &position, pool, baseCurrency.Denom)
		if err != nil {
			// Add log about error or not closed
			takeProfitLog = append(takeProfitLog, fmt.Sprintf("Position: Address:%s Id:%d cannot be liquidated due to err: %s", position.Address, position.Id, err.Error()))
		} else {
			/* *************************************************************************** */
			/* Start of kwak-indexer node implementation*/
			initialValue := math.LegacyNewDecFromInt(position.Collateral)

			// For perpetual positions, we calculate final value based on custody and liabilities
			finalValue := math.LegacyNewDecFromInt(position.Custody.Sub(position.Liabilities))
			profitLoss := finalValue.Sub(initialValue)
			var profitLossPerc math.LegacyDec
			if !initialValue.IsZero() {
				profitLossPerc = profitLoss.Quo(initialValue).Mul(math.LegacyNewDec(100))
			}

			indexer.QueueEvent(ctx, "/elys-event/perpetual/take-profit", indexerPerpetualTypes.TakeProfitEvent{
				Address:  position.Address,
				ID:       position.Id,
				Position: position.Position.String(),
				Collateral: indexerTypes.Token{
					Amount: position.Collateral.String(),
					Denom:  position.CollateralAsset,
				},
				Custody: indexerTypes.Token{
					Amount: position.Custody.String(),
					Denom:  position.CustodyAsset,
				},
				Liabilities: indexerTypes.Token{
					Amount: position.Liabilities.String(),
					Denom:  position.LiabilitiesAsset,
				},
				TakeProfitPrice:        position.TakeProfitPrice.String(),
				TakeProfitLiabilities:  position.TakeProfitLiabilities.String(),
				TakeProfitCustody:      position.TakeProfitCustody.String(),
				TakeProfitBorrowFactor: position.TakeProfitBorrowFactor.String(),
				OpenPrice:              position.OpenPrice.String(),
				Health:                 position.MtpHealth.String(),
				ProfitLoss:             profitLoss.String(),
				ProfitLossPerc:         profitLossPerc.String(),
			}, []string{position.Address})
			/* End of kwak-indexer node implementation*/
			/* *************************************************************************** */

			ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClosePositions,
				sdk.NewAttribute("address", position.Address),
				sdk.NewAttribute("id", strconv.FormatUint(position.Id, 10)),
			))
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClosePositions,
		sdk.NewAttribute("liquidations", strings.Join(liqLog, "\n")),
		sdk.NewAttribute("stop_loss", strings.Join(closeLog, "\n")),
		sdk.NewAttribute("take_profit", strings.Join(takeProfitLog, "\n")),
	))

	return &types.MsgClosePositionsResponse{}, nil
}
