package perpetual

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgUpdateTakeProfitPrice struct {
	Creator      string      `json:"creator"`
	ID           uint64      `json:"id"`
	Price        string      `json:"price"`
	Position     string      `json:"position"`
	Collateral   types.Token `json:"collateral"`
	OpenPrice    string      `json:"open_price"`
	CurrentPrice string      `json:"current_price"`
}

func (m MsgUpdateTakeProfitPrice) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
