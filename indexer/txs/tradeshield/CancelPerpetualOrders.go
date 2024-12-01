package tradeshield

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgCancelPerpetualOrders struct {
	OwnerAddress string   `json:"owner_address"`
	OrderIds     []uint64 `json:"order_ids"`
}

func (m MsgCancelPerpetualOrders) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
