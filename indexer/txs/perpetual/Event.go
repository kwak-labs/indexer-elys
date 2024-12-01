package perpetual

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type ClosePositionEvent struct {
	Address        string      `json:"address"`
	ID             uint64      `json:"id"`
	Collateral     types.Token `json:"collateral"`
	Custody        types.Token `json:"custody"`
	Liabilities    types.Token `json:"liabilities"`
	Health         string      `json:"health"`
	InitialValue   string      `json:"initial_value"`
	FinalValue     string      `json:"final_value"`
	ProfitLoss     string      `json:"profit_loss"`
	ProfitLossPerc string      `json:"profit_loss_perc"`
	OpenPrice      string      `json:"open_price"`
	Position       string      `json:"position"` // Long or Short
}

type TakeProfitEvent struct {
	Address                string      `json:"address"`
	ID                     uint64      `json:"id"`
	Position               string      `json:"position"` // Long or Short
	Collateral             types.Token `json:"collateral"`
	Custody                types.Token `json:"custody"`
	Liabilities            types.Token `json:"liabilities"`
	TakeProfitPrice        string      `json:"take_profit_price"`
	TakeProfitLiabilities  string      `json:"take_profit_liabilities"`
	TakeProfitCustody      string      `json:"take_profit_custody"`
	TakeProfitBorrowFactor string      `json:"take_profit_borrow_factor"`
	OpenPrice              string      `json:"open_price"`
	Health                 string      `json:"health"`
	ProfitLoss             string      `json:"profit_loss"`
	ProfitLossPerc         string      `json:"profit_loss_perc"`
}

type StopLossEvent struct {
	Address        string      `json:"address"`
	ID             uint64      `json:"id"`
	Position       string      `json:"position"`
	Collateral     types.Token `json:"collateral"`
	Custody        types.Token `json:"custody"`
	Liabilities    types.Token `json:"liabilities"`
	StopLossPrice  string      `json:"stop_loss_price"`
	OpenPrice      string      `json:"open_price"`
	Health         string      `json:"health"`
	ProfitLoss     string      `json:"profit_loss"`
	ProfitLossPerc string      `json:"profit_loss_perc"`
}

type UpdateStopLossEvent struct {
	ID       uint64 `json:"id"`
	Address  string `json:"address"`
	StopLoss string `json:"stop_loss"`
}

func (e ClosePositionEvent) Process(database types.DatabaseManager, event types.BaseEvent) (types.Response, error) {
	mergedData := types.GenericEvent{
		BaseEvent: event,
		Data:      e,
	}

	err := database.ProcessNewEvent(mergedData, event.Author)
	if err != nil {
		return types.Response{}, fmt.Errorf("error processing transaction: %w", err)
	}

	return types.Response{}, nil
}

func (e TakeProfitEvent) Process(database types.DatabaseManager, event types.BaseEvent) (types.Response, error) {
	mergedData := types.GenericEvent{
		BaseEvent: event,
		Data:      e,
	}

	err := database.ProcessNewEvent(mergedData, event.Author)
	if err != nil {
		return types.Response{}, fmt.Errorf("error processing transaction: %w", err)
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
		return types.Response{}, fmt.Errorf("error processing transaction: %w", err)
	}

	return types.Response{}, nil
}

func (e UpdateStopLossEvent) Process(database types.DatabaseManager, event types.BaseEvent) (types.Response, error) {
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
