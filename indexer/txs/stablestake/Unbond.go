package stablestake

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgUnbond struct {
	Creator         string      `json:"creator"`
	Amount          string      `json:"amount"`
	ShareDenom      string      `json:"share_denom"`
	RedemptionRate  string      `json:"redemption_rate"`
	RedemptionToken types.Token `json:"redemption_token"`
}

func (m MsgUnbond) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
