package leveragelp

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgClose struct {
	Creator     string   `json:"creator"`
	ID          uint64   `json:"id"`
	LpAmount    string   `json:"lp_amount"`
	RepayAmount string   `json:"repay_amount"`
	Position    Position `json:"position"`
}

type Position struct {
	ID             uint64      `json:"id"`
	Address        string      `json:"address"`
	Collateral     types.Token `json:"collateral"`
	Liabilities    string      `json:"liabilities"`
	PositionHealth string      `json:"position_health"`
}

func (m MsgClose) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
