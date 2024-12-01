package tradeshield

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type OrderExecutionLog struct {
	OrderID uint64 `json:"order_id"`
	Error   string `json:"error,omitempty"`
}

type MsgExecuteOrders struct {
	Creator           string              `json:"creator"`
	SpotOrderIds      []uint64            `json:"spot_order_ids"`
	PerpetualOrderIds []uint64            `json:"perpetual_order_ids"`
	SpotLogs          []OrderExecutionLog `json:"spot_logs"`
	PerpetualLogs     []OrderExecutionLog `json:"perpetual_logs"`
}

func (m MsgExecuteOrders) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
