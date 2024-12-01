package leveragelp

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type PositionRequest struct {
	Address string `json:"address"`
	ID      uint64 `json:"id"`
}

type MsgClosePositions struct {
	Creator    string            `json:"creator"`
	Liquidate  []PositionRequest `json:"liquidate"`
	StopLoss   []PositionRequest `json:"stop_loss"`
	LiquidLogs []string          `json:"liquid_logs"`
	CloseLogs  []string          `json:"close_logs"`
}

func (m MsgClosePositions) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
