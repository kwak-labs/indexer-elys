package tokenomics

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type Lockup struct {
	Amount          string `json:"amount"`
	UnlockTimestamp uint64 `json:"unlock_timestamp"`
}

type CommittedTokens struct {
	Denom   string   `json:"denom"`
	Amount  string   `json:"amount"`
	Lockups []Lockup `json:"lockups"`
}

type VestingTokens struct {
	Denom                string `json:"denom"`
	TotalAmount          string `json:"total_amount"`
	ClaimedAmount        string `json:"claimed_amount"`
	NumBlocks            int64  `json:"num_blocks"`
	StartBlock           int64  `json:"start_block"`
	VestStartedTimestamp int64  `json:"vest_started_timestamp"`
}

type MsgClaimAirdrop struct {
	Sender          string             `json:"sender"`
	AmountClaimed   types.Token        `json:"amount_claimed"`
	CommittedTokens []*CommittedTokens `json:"committed_tokens"`
	CommitsClaimed  types.Token        `json:"commits_claimed"`
	VestingTokens   []*VestingTokens   `json:"vesting_tokens"`
}

func (m MsgClaimAirdrop) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
