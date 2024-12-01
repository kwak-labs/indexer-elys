package common

type BaseEvent struct {
	EventID   uint64 `json:"event_id"`
	EventName string `json:"event_name"`
	EventTime uint64 `json:"event_time"`
}

type BaseTransaction struct {
	TransactionID uint64 `json:"transaction_id"`
	Author        string `json:"author"`
	Timestamp     uint64 `json:"timestamp"`
}
