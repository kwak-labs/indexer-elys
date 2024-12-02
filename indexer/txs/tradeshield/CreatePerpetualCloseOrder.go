package tradeshield

import (
	"fmt"

	"github.com/elys-network/elys/indexer/txs/tradeshield/common"
	"github.com/elys-network/elys/indexer/types"
)

type MsgCreatePerpetualCloseOrder struct {
	OwnerAddress string              `json:"owner_address"`
	TriggerPrice common.TriggerPrice `json:"trigger_price"`
	PositionID   uint64              `json:"position_id"`
	OrderID      uint64              `json:"order_id"`
}

func (m MsgCreatePerpetualCloseOrder) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
