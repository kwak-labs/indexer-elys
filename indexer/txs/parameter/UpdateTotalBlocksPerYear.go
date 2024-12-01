package parameter

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgUpdateTotalBlocksPerYear struct {
	Creator            string `json:"creator"`
	TotalBlocksPerYear uint64 `json:"total_blocks_per_year"`
}

func (m MsgUpdateTotalBlocksPerYear) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
