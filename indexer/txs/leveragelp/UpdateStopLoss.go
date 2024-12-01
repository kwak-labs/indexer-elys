package leveragelp

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type PositionStopLoss struct {
	ID          uint64      `json:"id"`
	Address     string      `json:"address"`
	Collateral  types.Token `json:"collateral"`
	Liabilities string      `json:"liabilities"`
	Health      string      `json:"health"`
	StopLoss    string      `json:"stop_loss"`
}

type MsgUpdateStopLoss struct {
	Creator   string           `json:"creator"`
	Position  uint64           `json:"position"`
	Price     string           `json:"price"`
	PoolID    uint64           `json:"pool_id"`
	Position_ PositionStopLoss `json:"position_details"`
}

func (m MsgUpdateStopLoss) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
