package commitments

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgUncommitTokens struct {
	Creator       string      `json:"creator"`
	Token         types.Token `json:"token"`
	IsLiquidation bool        `json:"is_liquidation"`
}

func (m MsgUncommitTokens) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
