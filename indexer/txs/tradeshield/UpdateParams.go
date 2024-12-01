package tradeshield

import (
	"fmt"

	common "github.com/elys-network/elys/indexer/txs/tradeshield/common"
	"github.com/elys-network/elys/indexer/types"
)

type MsgUpdateParams struct {
	Authority string        `json:"authority"`
	Params    common.Params `json:"params"`
}

func (m MsgUpdateParams) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
