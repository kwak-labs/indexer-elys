package leveragelp

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type AddPool struct {
	AmmPoolID   uint64 `json:"amm_pool_id"`
	LeverageMax string `json:"leverage_max"`
	Leverage    string `json:"leverage"`
}

type MsgAddPool struct {
	Authority string  `json:"authority"`
	Pool      AddPool `json:"pool"`
}

func (m MsgAddPool) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
