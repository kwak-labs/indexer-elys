package commitments

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgVestNow struct {
	Creator      string      `json:"creator"`
	Amount       string      `json:"amount"`
	Denom        string      `json:"denom"`
	VestAmount   types.Token `json:"vest_amount"`
	VestingDenom string      `json:"vesting_denom"`
}

func (m MsgVestNow) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
