package keeper

import (
	"context"
	"fmt"
	"strings"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerTradeshieldTypes "github.com/elys-network/elys/indexer/txs/tradeshield"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (k msgServer) ExecuteOrders(goCtx context.Context, msg *types.MsgExecuteOrders) (*types.MsgExecuteOrdersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	spotExecutionLogs := []indexerTradeshieldTypes.OrderExecutionLog{}
	perpExecutionLogs := []indexerTradeshieldTypes.OrderExecutionLog{}
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	spotLog := []string{}
	// loop through the spot orders and execute them
	for _, spotOrderId := range msg.SpotOrderIds {
		// get the spot order
		spotOrder, found := k.GetPendingSpotOrder(ctx, spotOrderId)
		if !found {
			return nil, types.ErrSpotOrderNotFound
		}

		var err error
		var res *ammtypes.MsgSwapByDenomResponse

		// dispatch based on the order type
		switch spotOrder.OrderType {
		case types.SpotOrderType_STOPLOSS:
			res, err = k.ExecuteStopLossOrder(ctx, spotOrder)
		case types.SpotOrderType_LIMITSELL:
			res, err = k.ExecuteLimitSellOrder(ctx, spotOrder)
		case types.SpotOrderType_LIMITBUY:
			res, err = k.ExecuteLimitBuyOrder(ctx, spotOrder)
		case types.SpotOrderType_MARKETBUY:
			res, err = k.ExecuteMarketBuyOrder(ctx, spotOrder)
		}

		if err != nil {
			errLog := fmt.Sprintf("Spot order Id:%d cannot be executed due to err: %s", spotOrderId, err.Error())
			spotLog = append(spotLog, errLog)
			/* *************************************************************************** */
			/* Start of kwak-indexer node implementation*/
			spotExecutionLogs = append(spotExecutionLogs, indexerTradeshieldTypes.OrderExecutionLog{
				OrderID: spotOrderId,
				Error:   err.Error(),
			})
			/* End of kwak-indexer node implementation*/
			/* *************************************************************************** */
		} else {
			ctx.EventManager().EmitEvent(types.NewExecuteSpotOrderEvt(spotOrder, res))
			/* *************************************************************************** */
			/* Start of kwak-indexer node implementation*/
			spotExecutionLogs = append(spotExecutionLogs, indexerTradeshieldTypes.OrderExecutionLog{
				OrderID: spotOrderId,
			})
			/* End of kwak-indexer node implementation*/
			/* *************************************************************************** */
		}
	}

	perpLog := []string{}
	// loop through the perpetual orders and execute them
	for _, perpetualOrderId := range msg.PerpetualOrderIds {
		perpetualOrder, found := k.GetPendingPerpetualOrder(ctx, perpetualOrderId)
		if !found {
			return nil, types.ErrPerpetualOrderNotFound
		}

		var err error

		switch perpetualOrder.PerpetualOrderType {
		case types.PerpetualOrderType_LIMITOPEN:
			err = k.ExecuteLimitOpenOrder(ctx, perpetualOrder)
		}

		if err != nil {
			errLog := fmt.Sprintf("Perpetual order Id:%d cannot be executed due to err: %s", perpetualOrderId, err.Error())
			perpLog = append(perpLog, errLog)
			/* *************************************************************************** */
			/* Start of kwak-indexer node implementation*/
			perpExecutionLogs = append(perpExecutionLogs, indexerTradeshieldTypes.OrderExecutionLog{
				OrderID: perpetualOrderId,
				Error:   err.Error(),
			})
			/* End of kwak-indexer node implementation*/
			/* *************************************************************************** */
		} else {
			/* *************************************************************************** */
			/* Start of kwak-indexer node implementation*/
			perpExecutionLogs = append(perpExecutionLogs, indexerTradeshieldTypes.OrderExecutionLog{
				OrderID: perpetualOrderId,
			})
			/* End of kwak-indexer node implementation*/
			/* *************************************************************************** */
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.TypeEvtExecuteOrders,
		sdk.NewAttribute("spot_orders", strings.Join(spotLog, "\n")),
		sdk.NewAttribute("perpetual_orders", strings.Join(perpLog, "\n")),
	))

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerTradeshieldTypes.MsgExecuteOrders{
		Creator:           msg.Creator,
		SpotOrderIds:      msg.SpotOrderIds,
		PerpetualOrderIds: msg.PerpetualOrderIds,
		SpotLogs:          spotExecutionLogs,
		PerpetualLogs:     perpExecutionLogs,
	}, []string{msg.Creator})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgExecuteOrdersResponse{}, nil
}
