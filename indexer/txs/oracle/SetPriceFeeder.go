package oracle

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgSetPriceFeeder struct {
	Feeder   string `json:"feeder"`
	IsActive bool   `json:"is_active"`
}

func (m MsgSetPriceFeeder) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
