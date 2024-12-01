// common/order_types.go

package common

import "github.com/elys-network/elys/indexer/types"

type OrderPrice struct {
	BaseDenom  string `json:"base_denom"`
	QuoteDenom string `json:"quote_denom"`
	Rate       string `json:"rate"`
}

type OrderDate struct {
	Height    uint64 `json:"height"`
	Timestamp uint64 `json:"timestamp"`
}

type BaseOrder struct {
	OrderID      uint64      `json:"order_id"`
	OwnerAddress string      `json:"owner_address"`
	OrderPrice   OrderPrice  `json:"order_price"`
	OrderAmount  types.Token `json:"order_amount"`
}

type OrderType int32

const (
	OrderType_UNKNOWN   OrderType = 0
	OrderType_SPOT      OrderType = 1
	OrderType_PERPETUAL OrderType = 2
)

type SpotOrder int32

const (
	SpotOrder_UNKNOWN    SpotOrder = 0
	SpotOrder_LIMIT_BUY  SpotOrder = 1
	SpotOrder_LIMIT_SELL SpotOrder = 2
	SpotOrder_MARKET_BUY SpotOrder = 3
)

type OrderStatus int32

const (
	Status_PENDING  OrderStatus = 0
	Status_EXECUTED OrderStatus = 1
	Status_CANCELED OrderStatus = 2
)
