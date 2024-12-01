package perpetual

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type Position int32

const (
	Position_UNSPECIFIED Position = 0
	Position_LONG        Position = 1
	Position_SHORT       Position = 2
)

type MsgOpen struct {
	Creator         string      `json:"creator"`
	Position        Position    `json:"position"`
	Leverage        string      `json:"leverage"`
	TradingAsset    string      `json:"trading_asset"`
	Collateral      types.Token `json:"collateral"`
	TakeProfitPrice string      `json:"take_profit_price"`
	StopLossPrice   string      `json:"stop_loss_price"`
	PoolID          uint64      `json:"pool_id"`
	PositionID      uint64      `json:"position_id"`
	OpenPrice       string      `json:"open_price"`
}

func (m MsgOpen) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
