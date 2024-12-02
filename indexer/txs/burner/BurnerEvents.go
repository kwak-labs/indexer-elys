package burner

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

// Event for tokens sent from zero address to module
type ZeroAddressTransferEvent struct {
	FromAddress string        `json:"from_address"`
	ToModule    string        `json:"to_module"`
	Coins       []types.Token `json:"coins"`
	Timestamp   string        `json:"timestamp"`
}

// Event for burning tokens
type TokenBurnEvent struct {
	Module    string        `json:"module"`
	Coins     []types.Token `json:"coins"`
	Timestamp string        `json:"timestamp"`
}

func (e ZeroAddressTransferEvent) Process(database types.DatabaseManager, event types.BaseEvent) (types.Response, error) {
	mergedData := types.GenericEvent{
		BaseEvent: event,
		Data:      e,
	}

	err := database.ProcessNewEvent(mergedData, event.Author)
	if err != nil {
		return types.Response{}, fmt.Errorf("error processing zero address transfer event: %w", err)
	}

	return types.Response{}, nil
}

func (e TokenBurnEvent) Process(database types.DatabaseManager, event types.BaseEvent) (types.Response, error) {
	mergedData := types.GenericEvent{
		BaseEvent: event,
		Data:      e,
	}

	err := database.ProcessNewEvent(mergedData, event.Author)
	if err != nil {
		return types.Response{}, fmt.Errorf("error processing token burn event: %w", err)
	}

	return types.Response{}, nil
}
