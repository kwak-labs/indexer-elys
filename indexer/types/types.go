package types

import (
	"time"
)

type FeeDetail struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

type Token struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

type BaseTransaction struct {
	BlockTime         time.Time   `json:"block_time"`
	Author            string      `json:"author"`
	IncludedAddresses []string    `json:"included_addresses"`
	BlockHeight       int64       `json:"block_height"`
	TxHash            string      `json:"tx_hash"`
	TxType            string      `json:"tx_type"`
	Fees              []FeeDetail `json:"fees"`
	GasLimit          string      `json:"gas_limit"`
	GasUsed           string      `json:"gas_used"`
	Memo              string      `json:"memo"`
	Status            string      `json:"status"`
}

type BaseEvent struct {
	BlockTime         time.Time `json:"block_time"`
	Author            string    `json:"author"`
	IncludedAddresses []string  `json:"included_addresses"`
	BlockHeight       int64     `json:"block_height"`
	EventType         string    `json:"event_type"`
}

type GenericTransaction struct {
	BaseTransaction BaseTransaction `json:"base_transaction"`
	Data            interface{}     `json:"data"`
}

type GenericEvent struct {
	BaseEvent BaseEvent   `json:"base_event"`
	Data      interface{} `json:"data"`
}

type GenericRecord struct {
	Transaction *GenericTransaction `json:"transaction,omitempty"`
	Event       *GenericEvent       `json:"event,omitempty"`
}

type Response struct{}

type DatabaseManager interface {
	ProcessNewTx(GenericTransaction, string) error
	ProcessNewEvent(GenericEvent, string) error
}

type Processor interface {
	Process(DatabaseManager, BaseTransaction) (Response, error)
}

type EventProcessor interface {
	Process(DatabaseManager, BaseEvent) (Response, error)
}

func (r GenericRecord) IsTransaction() bool {
	return r.Transaction != nil
}

func (r GenericRecord) IsEvent() bool {
	return r.Event != nil
}
