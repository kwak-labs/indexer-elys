package tokenomics

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgDeleteTimeBasedInflation struct {
	Authority        string `json:"authority"`
	StartBlockHeight uint64 `json:"start_block_height"`
	EndBlockHeight   uint64 `json:"end_block_height"`
}

func (m MsgDeleteTimeBasedInflation) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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