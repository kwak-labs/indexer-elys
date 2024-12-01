package tradeshield

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type TriggerPrice struct {
	TradingAssetDenom string `json:"trading_asset_denom"`
	Rate              string `json:"rate"`
}

type MsgCreatePerpetualOpenOrder struct {
	OwnerAddress    string       `json:"owner_address"`
	TriggerPrice    TriggerPrice `json:"trigger_price"`
	Collateral      types.Token  `json:"collateral"`
	TradingAsset    string       `json:"trading_asset"`
	Position        int32        `json:"position"`
	Leverage        string       `json:"leverage"`
	TakeProfitPrice string       `json:"take_profit_price"`
	StopLossPrice   string       `json:"stop_loss_price"`
	PoolID          uint64       `json:"pool_id"`
	OrderID         uint64       `json:"order_id"`
}

func (m MsgCreatePerpetualOpenOrder) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
