package tokenomics

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgDeleteAirdrop struct {
	Authority string `json:"authority"`
	Intent    string `json:"intent"`
}

func (m MsgDeleteAirdrop) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
