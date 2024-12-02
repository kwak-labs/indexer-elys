package oracle

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgRemoveAssetInfo struct {
	Authority string `json:"authority"`
	Denom     string `json:"denom"`
}

func (m MsgRemoveAssetInfo) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
