package indexer

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/elys-network/elys/indexer/txs/amm"
	"github.com/elys-network/elys/indexer/txs/assetprofile"
	"github.com/elys-network/elys/indexer/txs/burner"
	"github.com/elys-network/elys/indexer/txs/commitments"
	"github.com/elys-network/elys/indexer/txs/estaking"
	"github.com/elys-network/elys/indexer/txs/leveragelp"
	"github.com/elys-network/elys/indexer/txs/masterchef"
	"github.com/elys-network/elys/indexer/txs/oracle"
	"github.com/elys-network/elys/indexer/txs/parameter"
	"github.com/elys-network/elys/indexer/txs/perpetual"
	"github.com/elys-network/elys/indexer/txs/stablestake"
	"github.com/elys-network/elys/indexer/txs/tier"
	"github.com/elys-network/elys/indexer/txs/tokenomics"
	"github.com/elys-network/elys/indexer/txs/tradeshield"
	"github.com/elys-network/elys/indexer/types"
	indexerTypes "github.com/elys-network/elys/indexer/types"
)

var txRegistry = make(map[string]reflect.Type)
var eventRegistry = make(map[string]reflect.Type)

func init() {
	// Commitments
	RegisterTxType("/elys.commitment.MsgStake", reflect.TypeOf(commitments.MsgStake{}))
	RegisterTxType("/elys.commitment.MsgUnstake", reflect.TypeOf(commitments.MsgUnstake{}))
	RegisterTxType("/elys.commitment.MsgVestLiquid", reflect.TypeOf(commitments.MsgVestLiquid{}))
	RegisterTxType("/elys.commitment.MsgCancelVest", reflect.TypeOf(commitments.MsgCancelVest{}))
	RegisterTxType("/elys.commitment.MsgClaimVesting", reflect.TypeOf(commitments.MsgClaimVesting{}))
	RegisterTxType("/elys.commitment.MsgCommitClaimedRewards", reflect.TypeOf(commitments.MsgCommitClaimedRewards{}))
	RegisterTxType("/elys.commitment.MsgUncommitTokens", reflect.TypeOf(commitments.MsgUncommitTokens{}))
	RegisterTxType("/elys.commitment.MsgVest", reflect.TypeOf(commitments.MsgVest{}))
	RegisterTxType("/elys.commitment.MsgVestNow", reflect.TypeOf(commitments.MsgVestNow{}))
	RegisterTxType("/elys.commitment.MsgUpdateVestingInfo", reflect.TypeOf(commitments.MsgUpdateVestingInfo{}))

	// AMM
	RegisterTxType("/elys.amm.MsgCreatePool", reflect.TypeOf(amm.MsgCreatePool{}))
	RegisterTxType("/elys.amm.MsgJoinPool", reflect.TypeOf(amm.MsgJoinPool{}))
	RegisterTxType("/elys.amm.MsgExitPool", reflect.TypeOf(amm.MsgExitPool{}))
	RegisterTxType("/elys.amm.MsgSwapExactAmountIn", reflect.TypeOf(amm.MsgSwapExactAmountIn{}))
	RegisterTxType("/elys.amm.MsgSwapExactAmountOut", reflect.TypeOf(amm.MsgSwapExactAmountOut{}))
	RegisterTxType("/elys.amm.MsgSwapByDenom", reflect.TypeOf(amm.MsgSwapByDenom{}))
	RegisterTxType("/elys.amm.MsgUpdateParams", reflect.TypeOf(amm.MsgUpdateParams{}))
	RegisterTxType("/elys.amm.MsgUpdatePoolParams", reflect.TypeOf(amm.MsgUpdatePoolParams{}))
	RegisterTxType("/elys.amm.MsgFeedMultipleExternalLiquidity", reflect.TypeOf(amm.MsgFeedMultipleExternalLiquidity{}))

	// Perpetual
	RegisterTxType("/elys.perpetual.MsgOpen", reflect.TypeOf(perpetual.MsgOpen{}))
	RegisterTxType("/elys.perpetual.MsgClose", reflect.TypeOf(perpetual.MsgClose{}))
	RegisterTxType("/elys.perpetual.MsgUpdateStopLoss", reflect.TypeOf(perpetual.MsgUpdateStopLoss{}))
	RegisterTxType("/elys.perpetual.MsgClosePositions", reflect.TypeOf(perpetual.MsgClosePositions{}))
	RegisterTxType("/elys.perpetual.MsgUpdateParams", reflect.TypeOf(perpetual.MsgUpdateParams{}))
	RegisterTxType("/elys.perpetual.MsgWhitelist", reflect.TypeOf(perpetual.MsgWhitelist{}))
	RegisterTxType("/elys.perpetual.MsgDewhitelist", reflect.TypeOf(perpetual.MsgDewhitelist{}))
	RegisterTxType("/elys.perpetual.MsgUpdateTakeProfitPrice", reflect.TypeOf(perpetual.MsgUpdateTakeProfitPrice{}))

	// LeverageLP
	RegisterTxType("/elys.leveragelp.MsgOpen", reflect.TypeOf(leveragelp.MsgOpen{}))
	RegisterTxType("/elys.leveragelp.MsgClose", reflect.TypeOf(leveragelp.MsgClose{}))
	RegisterTxType("/elys.leveragelp.MsgClosePosition", reflect.TypeOf(leveragelp.MsgClosePositions{}))
	RegisterTxType("/elys.leveragelp.MsgClaimRewards", reflect.TypeOf(leveragelp.MsgClaimRewards{}))
	RegisterTxType("/elys.leveragelp.MsgUpdateStopLoss", reflect.TypeOf(leveragelp.MsgUpdateStopLoss{}))
	RegisterTxType("/elys.leveragelp.MsgAddPool", reflect.TypeOf(leveragelp.MsgAddPool{}))
	RegisterTxType("/elys.leveragelp.MsgUpdateParams", reflect.TypeOf(leveragelp.MsgUpdateParams{}))
	RegisterTxType("/elys.leveragelp.MsgWhitelist", reflect.TypeOf(leveragelp.MsgWhitelist{}))
	RegisterTxType("/elys.leveragelp.MsgDewhitelist", reflect.TypeOf(leveragelp.MsgDewhitelist{}))
	RegisterTxType("/elys.leveragelp.MsgRemovePool", reflect.TypeOf(leveragelp.MsgRemovePool{}))

	// Oracle
	RegisterTxType("/elys.oracle.MsgFeedPrice", reflect.TypeOf(oracle.MsgFeedPrice{}))
	RegisterTxType("/elys.oracle.MsgFeedMultiplePrices", reflect.TypeOf(oracle.MsgFeedMultiplePrices{}))
	RegisterTxType("/elys.oracle.MsgSetPriceFeeder", reflect.TypeOf(oracle.MsgSetPriceFeeder{}))
	RegisterTxType("/elys.oracle.MsgDeletePriceFeeder", reflect.TypeOf(oracle.MsgDeletePriceFeeder{}))
	RegisterTxType("/elys.oracle.MsgCreateAssetInfo", reflect.TypeOf(oracle.MsgCreateAssetInfo{}))
	RegisterTxType("/elys.oracle.MsgRemoveAssetInfo", reflect.TypeOf(oracle.MsgRemoveAssetInfo{}))
	RegisterTxType("/elys.oracle.MsgAddPriceFeeders", reflect.TypeOf(oracle.MsgAddPriceFeeders{}))
	RegisterTxType("/elys.oracle.MsgRemovePriceFeeders", reflect.TypeOf(oracle.MsgRemovePriceFeeders{}))
	RegisterTxType("/elys.oracle.MsgUpdateParams", reflect.TypeOf(oracle.MsgUpdateParams{}))

	// Parameter
	RegisterTxType("/elys.parameter.MsgUpdateMinCommission", reflect.TypeOf(parameter.MsgUpdateMinCommission{}))
	RegisterTxType("/elys.parameter.MsgUpdateMaxVotingPower", reflect.TypeOf(parameter.MsgUpdateMaxVotingPower{}))
	RegisterTxType("/elys.parameter.MsgUpdateMinSelfDelegation", reflect.TypeOf(parameter.MsgUpdateMinSelfDelegation{}))
	RegisterTxType("/elys.parameter.MsgUpdateTotalBlocksPerYear", reflect.TypeOf(parameter.MsgUpdateTotalBlocksPerYear{}))
	RegisterTxType("/elys.parameter.MsgUpdateRewardsDataLifetime", reflect.TypeOf(parameter.MsgUpdateRewardsDataLifetime{}))

	// StableStake
	RegisterTxType("/elys.stablestake.MsgBond", reflect.TypeOf(stablestake.MsgBond{}))
	RegisterTxType("/elys.stablestake.MsgUnbond", reflect.TypeOf(stablestake.MsgUnbond{}))
	RegisterTxType("/elys.stablestake.MsgUpdateParams", reflect.TypeOf(stablestake.MsgUpdateParams{}))

	// TradeShield
	RegisterTxType("/elys.tradeshield.MsgCreateSpotOrder", reflect.TypeOf(tradeshield.MsgCreateSpotOrder{}))
	RegisterTxType("/elys.tradeshield.MsgCancelSpotOrders", reflect.TypeOf(tradeshield.MsgCancelSpotOrders{}))
	RegisterTxType("/elys.tradeshield.MsgCreatePerpetualOrder", reflect.TypeOf(tradeshield.MsgCreatePerpetualOpenOrder{}))
	RegisterTxType("/elys.tradeshield.MsgCancelPerpetualOrder", reflect.TypeOf(tradeshield.MsgCancelPerpetualOrder{}))
	RegisterTxType("/elys.tradeshield.MsgCancelPerpetualOrders", reflect.TypeOf(tradeshield.MsgCancelPerpetualOrders{}))
	RegisterTxType("/elys.tradeshield.MsgUpdatePerpetualOrder", reflect.TypeOf(tradeshield.MsgUpdatePerpetualOrder{}))
	RegisterTxType("/elys.tradeshield.MsgExecuteOrders", reflect.TypeOf(tradeshield.MsgExecuteOrders{}))
	RegisterTxType("/elys.tradeshield.MsgUpdateParams", reflect.TypeOf(tradeshield.MsgUpdateParams{}))
	RegisterTxType("/elys.tradeshield.MsgUpdateSpotOrder", reflect.TypeOf(tradeshield.MsgUpdateSpotOrder{}))
	RegisterTxType("/elys.tradeshield.MsgCancelSpotOrder", reflect.TypeOf(tradeshield.MsgCancelSpotOrder{}))
	RegisterTxType("/elys.tradeshield.MsgCreatePerpetualCloseOrder", reflect.TypeOf(tradeshield.MsgCreatePerpetualCloseOrder{}))

	// Asset Profile
	RegisterTxType("/elys.assetprofile.MsgAddEntry", reflect.TypeOf(assetprofile.MsgAddEntry{}))
	RegisterTxType("/elys.assetprofile.MsgUpdateEntry", reflect.TypeOf(assetprofile.MsgUpdateEntry{}))
	RegisterTxType("/elys.assetprofile.MsgDeleteEntry", reflect.TypeOf(assetprofile.MsgDeleteEntry{}))

	// Masterchef
	RegisterTxType("/elys.masterchef.MsgClaimRewards", reflect.TypeOf(masterchef.MsgClaimRewards{}))
	RegisterTxType("/elys.masterchef.MsgAddExternalRewardDenom", reflect.TypeOf(masterchef.MsgAddExternalRewardDenom{}))
	RegisterTxType("/elys.masterchef.MsgAddExternalIncentive", reflect.TypeOf(masterchef.MsgAddExternalIncentive{}))
	RegisterTxType("/elys.masterchef.MsgUpdateParams", reflect.TypeOf(masterchef.MsgUpdateParams{}))
	RegisterTxType("/elys.masterchef.MsgUpdatePoolMultipliers", reflect.TypeOf(masterchef.MsgUpdatePoolMultipliers{}))
	RegisterTxType("/elys.masterchef.MsgTogglePoolEdenRewards", reflect.TypeOf(masterchef.MsgTogglePoolEdenRewards{}))

	// EStakingz
	RegisterTxType("/elys.estaking.MsgUpdateParams", reflect.TypeOf(estaking.MsgUpdateParams{}))
	RegisterTxType("/elys.estaking.MsgWithdrawAllRewards", reflect.TypeOf(estaking.MsgWithdrawAllRewards{}))
	RegisterTxType("/elys.estaking.MsgWithdrawElysStakingRewards", reflect.TypeOf(estaking.MsgWithdrawElysStakingRewards{}))
	RegisterTxType("/elys.estaking.MsgWithdrawReward", reflect.TypeOf(estaking.MsgWithdrawReward{}))

	// Burner
	RegisterTxType("/elys.burner.MsgUpdateParams", reflect.TypeOf(burner.MsgUpdateParams{}))

	// Tier
	RegisterTxType("/elys.tier.MsgSetPortfolio", reflect.TypeOf(tier.MsgSetPortfolio{}))

	// Tokenomics
	RegisterTxType("/elys.tokenomics.MsgCreateAirdrop", reflect.TypeOf(tokenomics.MsgCreateAirdrop{}))
	RegisterTxType("/elys.tokenomics.MsgUpdateAirdrop", reflect.TypeOf(tokenomics.MsgUpdateAirdrop{}))
	RegisterTxType("/elys.tokenomics.MsgDeleteAirdrop", reflect.TypeOf(tokenomics.MsgDeleteAirdrop{}))
	RegisterTxType("/elys.tokenomics.MsgClaimAirdrop", reflect.TypeOf(tokenomics.MsgClaimAirdrop{}))
	RegisterTxType("/elys.tokenomics.MsgUpdateGenesisInflation", reflect.TypeOf(tokenomics.MsgUpdateGenesisInflation{}))
	RegisterTxType("/elys.tokenomics.MsgCreateTimeBasedInflation", reflect.TypeOf(tokenomics.MsgCreateTimeBasedInflation{}))
	RegisterTxType("/elys.tokenomics.MsgUpdateTimeBasedInflation", reflect.TypeOf(tokenomics.MsgUpdateTimeBasedInflation{}))
	RegisterTxType("/elys.tokenomics.MsgDeleteTimeBasedInflation", reflect.TypeOf(tokenomics.MsgDeleteTimeBasedInflation{}))

	// Register Events
	RegisterEventType("/elys-event/leveragelp/liquidation", reflect.TypeOf(leveragelp.LiquidationEvent{}))
	RegisterEventType("/elys-event/leveragelp/stop-loss", reflect.TypeOf(leveragelp.StopLossEvent{}))
	RegisterEventType("/elys-event/masterchef/claim-rewards", reflect.TypeOf(masterchef.ClaimRewardsEvent{}))
	RegisterEventType("/elys-event/perpetual/liquidation", reflect.TypeOf(perpetual.LiquidationEvent{}))
	RegisterEventType("/elys-event/perpetual/stop-loss", reflect.TypeOf(perpetual.StopLossEvent{}))
	RegisterEventType("/elys-event/perpetual/take-profit", reflect.TypeOf(perpetual.TakeProfitEvent{}))
	RegisterEventType("/elys-event/tradeshield/stop-loss", reflect.TypeOf(tradeshield.StopLossExecutionEvent{}))
	RegisterEventType("/elys-event/tradeshield/limit-sell", reflect.TypeOf(tradeshield.LimitSellExecutionEvent{}))
	RegisterEventType("/elys-event/tradeshield/limit-buy", reflect.TypeOf(tradeshield.LimitOrderExecutionEvent{}))
	RegisterEventType("/elys-event/tradeshield/market-buy", reflect.TypeOf(tradeshield.MarketOrderExecutionEvent{}))
}

func RegisterTxType(txType string, dataType reflect.Type) {
	txRegistry[txType] = dataType
}

func RegisterEventType(eventType string, dataType reflect.Type) {
	eventRegistry[eventType] = dataType
}

func ParseTransaction(tx types.GenericTransaction) (string, types.Processor, error) {
	txType := tx.BaseTransaction.TxType

	dataType, ok := txRegistry[txType]
	if !ok {
		return "", nil, fmt.Errorf("unknown transaction type: %s", txType)
	}

	dataValue := reflect.New(dataType).Interface()
	dataBytes, err := json.Marshal(tx.Data)
	if err != nil {
		return "", nil, fmt.Errorf("error marshaling data: %w", err)
	}

	err = json.Unmarshal(dataBytes, dataValue)
	if err != nil {
		return "", nil, fmt.Errorf("error unmarshaling to %s: %w", dataType.Name(), err)
	}

	processor, ok := reflect.ValueOf(dataValue).Elem().Interface().(types.Processor)
	if !ok {
		return "", nil, fmt.Errorf("type %s does not implement Processor", dataType.Name())
	}

	return txType, processor, nil
}

func ParseEvent(event types.GenericEvent) (string, types.EventProcessor, error) {
	eventType := event.BaseEvent.EventType

	dataType, ok := eventRegistry[eventType]
	if !ok {
		return "", nil, fmt.Errorf("unknown event type: %s", eventType)
	}

	dataValue := reflect.New(dataType).Interface()
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return "", nil, fmt.Errorf("error marshaling event data: %w", err)
	}

	err = json.Unmarshal(dataBytes, dataValue)
	if err != nil {
		return "", nil, fmt.Errorf("error unmarshaling to %s: %w", dataType.Name(), err)
	}

	processor, ok := reflect.ValueOf(dataValue).Elem().Interface().(types.EventProcessor)
	if !ok {
		return "", nil, fmt.Errorf("type %s does not implement EventProcessor", dataType.Name())
	}

	return eventType, processor, nil
}

func ParseRecord(record indexerTypes.GenericRecord) (string, interface{}, error) {
	if record.IsTransaction() {
		return ParseTransaction(*record.Transaction)
	} else if record.IsEvent() {
		return ParseEvent(*record.Event)
	}
	return "", nil, fmt.Errorf("record contains neither transaction nor event")
}
