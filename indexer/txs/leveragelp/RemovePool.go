package leveragelp

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgRemovePool struct {
	Authority string `json:"authority"`
	PoolID    uint64 `json:"pool_id"`
}

func (m MsgRemovePool) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
