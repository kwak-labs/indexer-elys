package perpetual

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgClosePositions struct {
	Creator    string           `json:"creator"`
	Liquidate  []PositionResult `json:"liquidate"`
	StopLoss   []PositionResult `json:"stop_loss"`
	TakeProfit []PositionResult `json:"take_profit"`
}

type PositionResult struct {
	Address        string      `json:"address"`
	ID             uint64      `json:"id"`
	Position       string      `json:"position"`
	Collateral     types.Token `json:"collateral"`
	Custody        types.Token `json:"custody"`
	Liabilities    types.Token `json:"liabilities"`
	Health         string      `json:"health"`
	InitialValue   string      `json:"initial_value"`
	FinalValue     string      `json:"final_value"`
	ProfitLoss     string      `json:"profit_loss"`
	ProfitLossPerc string      `json:"profit_loss_perc"`
	OpenPrice      string      `json:"open_price"`

	// Optional fields for TakeProfit
	TakeProfitPrice        string `json:"take_profit_price,omitempty"`
	TakeProfitLiabilities  string `json:"take_profit_liabilities,omitempty"`
	TakeProfitCustody      string `json:"take_profit_custody,omitempty"`
	TakeProfitBorrowFactor string `json:"take_profit_borrow_factor,omitempty"`

	// Optional fields for StopLoss
	StopLossPrice string `json:"stop_loss_price,omitempty"`
}

func (m MsgClosePositions) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
	mergedData := types.GenericTransaction{
		BaseTransaction: transaction,
		Data:            m,
	}

	err := database.ProcessNewTx(mergedData, transaction.Author)
	if err != nil {
		return types.Response{}, fmt.Errorf("error processing transaction: %w", err)
	}

	return types.Response{}, nil
}
