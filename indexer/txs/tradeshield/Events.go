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
	OrderTargetDenom string               `json:"order_target_denom,omitempty"`
	Status           common.OrderStatus   `json:"status"`
	Date             common.OrderDate     `json:"date,omitempty"`
	MarketPrice      string               `json:"market_price,omitempty"`
	TriggerPrice     *common.TriggerPrice `json:"trigger_price,omitempty"` // For stop-loss or perpetual
	SwapOutput       *types.Token         `json:"swap_output,omitempty"`   // For limit or spot
	ReceivedAmount   *types.Token         `json:"received_amount,omitempty"`
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
	OrderTargetDenom string             `json:"order_target_denom"`
	MarketPrice      string             `json:"market_price"`
	ExecutionStatus  common.OrderStatus `json:"execution_status"`
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
	OrderType      common.OrderType `json:"order_type"`
	Date           common.OrderDate `json:"date"`
	ReceivedAmount types.Token      `json:"received_amount"`
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
	OrderType        common.SpotOrder   `json:"order_type"`
	OrderTargetDenom string             `json:"order_target_denom"`
	Status           common.OrderStatus `json:"status"`
	Date             common.OrderDate   `json:"date"`
	SwapOutput       types.Token        `json:"swap_output"`
	MarketPrice      string             `json:"market_price"`
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
	SwapOutput  types.Token      `json:"swap_output"`
	MarketPrice string           `json:"market_price"`
	TargetDenom string           `json:"target_denom"`
	Date        common.OrderDate `json:"date"`
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
	OrderType   common.OrderType `json:"order_type"`
	Date        common.OrderDate `json:"date"`
	TargetDenom string           `json:"target_denom"`
	SwapOutput  types.Token      `json:"swap_output"`
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
