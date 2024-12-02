package tradeshield

import (
	"fmt"

	"github.com/elys-network/elys/indexer/txs/tradeshield/common"
	"github.com/elys-network/elys/indexer/types"
)

// Unified Execution Event
type ExecutionEvent struct {
	common.BaseOrder
	OrderType        common.OrderType     `json:"order_type"`
	OrderTargetDenom string               `json:"order_target_denom"`
	Status           common.OrderStatus   `json:"status"`
	Date             common.OrderDate     `json:"date"`
	MarketPrice      string               `json:"market_price"`
	TriggerPrice     *common.TriggerPrice `json:"trigger_price,omitempty"`
	SwapOutput       *types.Token         `json:"swap_output"`
	ReceivedAmount   *types.Token         `json:"received_amount"`
	SpotPrice        string               `json:"spot_price"`
	SwapFee          string               `json:"swap_fee"`
	Discount         string               `json:"discount"`
	Recipient        string               `json:"recipient"`
}

func (e ExecutionEvent) Process(database types.DatabaseManager, event types.BaseEvent) (types.Response, error) {
	mergedData := types.GenericEvent{
		BaseEvent: event,
		Data:      e,
	}

	err := database.ProcessNewEvent(mergedData, event.Author)
	if err != nil {
		return types.Response{}, fmt.Errorf("error processing execution event: %w", err)
	}

	return types.Response{}, nil
}

// LimitSellExecutionEvent represents the event for a limit sell execution
type LimitSellExecutionEvent struct {
	common.BaseOrder
	OrderType        common.OrderType   `json:"order_type"`
	OrderTargetDenom string             `json:"order_target_denom"`
	Status           common.OrderStatus `json:"status"`
	Date             common.OrderDate   `json:"date"`
	MarketPrice      string             `json:"market_price"`
	SwapOutput       types.Token        `json:"swap_output"`
	SpotPrice        string             `json:"spot_price"`
	SwapFee          string             `json:"swap_fee"`
	Discount         string             `json:"discount"`
	Recipient        string             `json:"recipient"`
}

func (e LimitSellExecutionEvent) Process(database types.DatabaseManager, event types.BaseEvent) (types.Response, error) {
	mergedData := types.GenericEvent{
		BaseEvent: event,
		Data:      e,
	}

	err := database.ProcessNewEvent(mergedData, event.Author)
	if err != nil {
		return types.Response{}, fmt.Errorf("error processing event: %w", err)
	}

	return types.Response{}, nil
}

// SpotOrderExecution represents the event for a spot order execution
type SpotOrderExecution struct {
	common.BaseOrder
	OrderType        common.OrderType   `json:"order_type"`
	OrderTargetDenom string             `json:"order_target_denom"`
	Status           common.OrderStatus `json:"status"`
	Date             common.OrderDate   `json:"date"`
	MarketPrice      string             `json:"market_price"`
	SwapOutput       types.Token        `json:"swap_output"`
	ReceivedAmount   types.Token        `json:"received_amount"`
	SpotPrice        string             `json:"spot_price"`
	SwapFee          string             `json:"swap_fee"`
	Discount         string             `json:"discount"`
	Recipient        string             `json:"recipient"`
}

func (e SpotOrderExecution) Process(database types.DatabaseManager, event types.BaseEvent) (types.Response, error) {
	mergedData := types.GenericEvent{
		BaseEvent: event,
		Data:      e,
	}

	err := database.ProcessNewEvent(mergedData, event.Author)
	if err != nil {
		return types.Response{}, fmt.Errorf("error processing event: %w", err)
	}

	return types.Response{}, nil
}

// LimitOrderExecutionEvent represents the event for a limit buy execution
type LimitOrderExecutionEvent struct {
	common.BaseOrder
	OrderType        common.OrderType   `json:"order_type"`
	OrderTargetDenom string             `json:"order_target_denom"`
	Status           common.OrderStatus `json:"status"`
	Date             common.OrderDate   `json:"date"`
	MarketPrice      string             `json:"market_price"`
	SwapOutput       types.Token        `json:"swap_output"`
	SpotPrice        string             `json:"spot_price"`
	SwapFee          string             `json:"swap_fee"`
	Discount         string             `json:"discount"`
	Recipient        string             `json:"recipient"`
}

func (e LimitOrderExecutionEvent) Process(database types.DatabaseManager, event types.BaseEvent) (types.Response, error) {
	mergedData := types.GenericEvent{
		BaseEvent: event,
		Data:      e,
	}

	err := database.ProcessNewEvent(mergedData, event.Author)
	if err != nil {
		return types.Response{}, fmt.Errorf("error processing event: %w", err)
	}

	return types.Response{}, nil
}

// StopLossExecutionEvent represents the event for a stop-loss execution
type StopLossExecutionEvent struct {
	common.BaseOrder
	OrderType        common.OrderType    `json:"order_type"`
	OrderTargetDenom string              `json:"order_target_denom"`
	Status           common.OrderStatus  `json:"status"`
	Date             common.OrderDate    `json:"date"`
	MarketPrice      string              `json:"market_price"`
	SwapOutput       types.Token         `json:"swap_output"`
	TriggerPrice     common.TriggerPrice `json:"trigger_price"`
	SpotPrice        string              `json:"spot_price"`
	SwapFee          string              `json:"swap_fee"`
	Discount         string              `json:"discount"`
	Recipient        string              `json:"recipient"`
}

func (e StopLossExecutionEvent) Process(database types.DatabaseManager, event types.BaseEvent) (types.Response, error) {
	mergedData := types.GenericEvent{
		BaseEvent: event,
		Data:      e,
	}

	err := database.ProcessNewEvent(mergedData, event.Author)
	if err != nil {
		return types.Response{}, fmt.Errorf("error processing event: %w", err)
	}

	return types.Response{}, nil
}

type MarketOrderExecutionEvent struct {
	common.BaseOrder
	OrderType        common.OrderType   `json:"order_type"`
	OrderTargetDenom string             `json:"order_target_denom"`
	Status           common.OrderStatus `json:"status"`
	Date             common.OrderDate   `json:"date"`
	MarketPrice      string             `json:"market_price"`
	SwapOutput       types.Token        `json:"swap_output"`
	SpotPrice        string             `json:"spot_price"`
	SwapFee          string             `json:"swap_fee"`
	Discount         string             `json:"discount"`
	Recipient        string             `json:"recipient"`
}

func (e MarketOrderExecutionEvent) Process(database types.DatabaseManager, event types.BaseEvent) (types.Response, error) {
	mergedData := types.GenericEvent{
		BaseEvent: event,
		Data:      e,
	}

	err := database.ProcessNewEvent(mergedData, event.Author)
	if err != nil {
		return types.Response{}, fmt.Errorf("error processing market order execution event: %w", err)
	}

	return types.Response{}, nil
}

// ExecutionLog represents the log for batch execution processing
type ExecutionLog struct {
	OrderID uint64 `json:"order_id"`
	Error   string `json:"error,omitempty"`
}
