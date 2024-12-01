package parameter

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgUpdateMinCommission struct {
	Creator       string `json:"creator"`
	MinCommission string `json:"min_commission"`
}

func (m MsgUpdateMinCommission) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
