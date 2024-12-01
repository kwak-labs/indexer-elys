package leveragelp

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type LiquidationEvent struct {
	PositionID     uint64      `json:"position_id"`
	Address        string      `json:"address"`
	Collateral     types.Token `json:"collateral"`
	RepayAmount    string      `json:"repay_amount"`
	Liabilities    string      `json:"liabilities"`
	Health         string      `json:"health"`
	InitialValue   string      `json:"initial_value"`    // Initial position value (collateral)
	FinalValue     string      `json:"final_value"`      // Final value (repayAmount - liabilities)
	ProfitLoss     string      `json:"profit_loss"`      // Actual P/L
	ProfitLossPerc string      `json:"profit_loss_perc"` // P/L as percentage
}

type StopLossEvent struct {
	PositionID      uint64      `json:"position_id"`
	Address         string      `json:"address"`
	Collateral      types.Token `json:"collateral"`
	RepayAmount     string      `json:"repay_amount"`
	Liabilities     string      `json:"liabilities"`
	Health          string      `json:"health"`
	StopLossPrice   string      `json:"stop_loss_price"`
	LpTokenPrice    string      `json:"lp_token_price"`
	InitialValue    string      `json:"initial_value"`    // Initial position value
	FinalValue      string      `json:"final_value"`      // Final value at stop loss
	ProfitLoss      string      `json:"profit_loss"`      // Calculated P&L
	ProfitLossPerc  string      `json:"profit_loss_perc"` // P&L as percentage
	RemainingAmount string      `json:"remaining_amount"` // Amount returned to user after closure
}

func (e LiquidationEvent) Process(database types.DatabaseManager, event types.BaseEvent) (types.Response, error) {
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

func (e StopLossEvent) Process(database types.DatabaseManager, event types.BaseEvent) (types.Response, error) {
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
