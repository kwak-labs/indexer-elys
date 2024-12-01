package masterchef

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type ClaimRewardsEvent struct {
	Sender      string        `json:"sender"`
	Recipient   string        `json:"recipient"`
	PoolIDs     []uint64      `json:"pool_ids"`
	RewardCoins []types.Token `json:"reward_coins"`
}

func (e ClaimRewardsEvent) Process(database types.DatabaseManager, event types.BaseEvent) (types.Response, error) {
	mergedData := types.GenericEvent{
		BaseEvent: event,
		Data:      e,
	}

	err := database.ProcessNewEvent(mergedData, event.Author)
	if err != nil {
		return types.Response{}, fmt.Errorf("error processing event: %w", err)
	}

	return types.Response{}, nil
}
