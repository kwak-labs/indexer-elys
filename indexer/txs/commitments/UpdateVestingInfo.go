package commitments

import (
	"fmt"

	"github.com/elys-network/elys/indexer/types"
)

type MsgUpdateVestingInfo struct {
	Authority      string `json:"authority"`
	BaseDenom      string `json:"base_denom"`
	VestingDenom   string `json:"vesting_denom"`
	NumBlocks      int64  `json:"num_blocks"`
	VestNowFactor  int64  `json:"vest_now_factor"`
	NumMaxVestings int64  `json:"num_max_vestings"`
}

func (m MsgUpdateVestingInfo) Process(database types.DatabaseManager, transaction types.BaseTransaction) (types.Response, error) {
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
