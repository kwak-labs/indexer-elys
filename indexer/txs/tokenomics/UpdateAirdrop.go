package tokenomics

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgUpdateAirdrop struct {
	Authority string `json:"authority"`
	Intent    string `json:"intent"`
	Amount    uint64 `json:"amount"`
	Expiry    uint64 `json:"expiry"`
}

func (m MsgUpdateAirdrop) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
