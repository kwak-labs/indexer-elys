package tradeshield

import (
	"fmt"

	"github.com/elys-network/elys/indexer/txs/tradeshield/common"
	"github.com/elys-network/elys/indexer/types"
	"google.golang.org/genproto/googleapis/type/decimal"
)

// Update MsgCreateSpotOrder to use BaseOrder
type MsgCreateSpotOrder struct {
	BaseOrder        common.BaseOrder `json:"base_order"`
	OrderType        common.OrderType `json:"order_type"`
	OrderTargetDenom string           `json:"order_target_denom"`
	StopPrice        *decimal.Decimal `json:"stop_price,omitempty"`
}

type MsgUpdateSpotOrder struct {
	OrderID      uint64            `json:"order_id"`
	OwnerAddress string            `json:"owner_address"`
	OrderPrice   common.OrderPrice `json:"order_price"`
	StopPrice    *decimal.Decimal  `json:"stop_price,omitempty"`
}

type MsgCancelSpotOrder struct {
	OwnerAddress string `json:"owner_address"`
	OrderId      uint64 `json:"order_id"`
}

type MsgCancelSpotOrders struct {
	Creator      string   `json:"creator"`
	SpotOrderIds []uint64 `json:"spot_order_ids"`
}

func (m MsgCreateSpotOrder) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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

func (m MsgUpdateSpotOrder) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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

func (m MsgCancelSpotOrder) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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

func (m MsgCancelSpotOrders) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
