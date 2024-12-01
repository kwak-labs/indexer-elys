package staking

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	abci "github.com/cometbft/cometbft/abci/types"
)

var (
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
	_ module.HasServices         = AppModule{}
	_ module.HasInvariants       = AppModule{}
	_ module.HasABCIGenesis      = AppModule{}
	_ module.HasABCIEndBlock     = AppModule{}

	_ appmodule.AppModule       = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
)

// AppModule embeds the Cosmos SDK's x/staking AppModuleBasic.
type AppModuleBasic struct {
	staking.AppModuleBasic
}

// AppModule embeds the Cosmos SDK's x/staking AppModule where we only override
// specific methods.
type AppModule struct {
	// embed the Cosmos SDK's x/staking AppModule
	staking.AppModule
	cdc        codec.Codec
	keeper     *keeper.Keeper
	accKeeper  types.AccountKeeper
	bankKeeper types.BankKeeper
}

// NewAppModule creates a new AppModule object using the native x/staking module
// AppModule constructor.
func NewAppModule(cdc codec.Codec, keeper *keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper, subspace exported.Subspace) AppModule {
	stakingAppMod := staking.NewAppModule(cdc, keeper, ak, bk, subspace)
	return AppModule{
		AppModule:  stakingAppMod,
		keeper:     keeper,
		accKeeper:  ak,
		bankKeeper: bk,
		cdc:        cdc,
	}
}

// InitGenesis delegates the InitGenesis call to the underlying x/staking module,
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState

	cdc.MustUnmarshalJSON(data, &genesisState)
	return am.keeper.InitGenesis(ctx, &genesisState)
}

// EndBlock delegates the EndBlock call to the underlying x/staking module,
func (am AppModule) EndBlock(goCtx context.Context) ([]abci.ValidatorUpdate, error) {
	sdkCtx := sdk.UnwrapSDKContext(goCtx)

	// Patch 1: skip UnbondAllMatureValidators
	if sdkCtx.BlockHeight() == 11517072 {
		return am.BlockValidatorUpdatesPatch1(goCtx)
	}

	// Patch 2: use custom UnbondAllMatureValidators
	if sdkCtx.BlockHeight() >= 11517073 {
		return am.BlockValidatorUpdatesPatch2(goCtx)
	}

	return am.keeper.BlockValidatorUpdates(goCtx)
}

// BeginBlock returns the begin blocker for the staking module.
func (am AppModule) BeginBlock(ctx context.Context) error {
	return am.keeper.BeginBlocker(ctx)
}

// Patch 1: skip UnbondAllMatureValidators
func (am AppModule) BlockValidatorUpdatesPatch1(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	validatorAddressCodec := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())
	AddressCodec := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	// Calculate validator set changes.
	//
	// NOTE: ApplyAndReturnValidatorSetUpdates has to come before
	// UnbondAllMatureValidatorQueue.
	// This fixes a bug when the unbonding period is instant (is the case in
	// some of the tests). The test expected the validator to be completely
	// unbonded after the Endblocker (go from Bonded -> Unbonding during
	// ApplyAndReturnValidatorSetUpdates and then Unbonding -> Unbonded during
	// UnbondAllMatureValidatorQueue).
	validatorUpdates, err := am.keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds, err := am.keeper.DequeueAllMatureUBDQueue(ctx, sdkCtx.BlockHeader().Time)
	if err != nil {
		return nil, err
	}

	for _, dvPair := range matureUnbonds {
		addr, err := validatorAddressCodec.StringToBytes(dvPair.ValidatorAddress)
		if err != nil {
			return nil, err
		}
		delegatorAddress, err := AddressCodec.StringToBytes(dvPair.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		balances, err := am.keeper.CompleteUnbonding(ctx, delegatorAddress, addr)
		if err != nil {
			continue
		}

		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbonding,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, dvPair.ValidatorAddress),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress),
			),
		)
	}

	// Remove all mature redelegations from the red queue.
	matureRedelegations, err := am.keeper.DequeueAllMatureRedelegationQueue(ctx, sdkCtx.BlockHeader().Time)
	if err != nil {
		return nil, err
	}

	for _, dvvTriplet := range matureRedelegations {
		valSrcAddr, err := validatorAddressCodec.StringToBytes(dvvTriplet.ValidatorSrcAddress)
		if err != nil {
			return nil, err
		}
		valDstAddr, err := validatorAddressCodec.StringToBytes(dvvTriplet.ValidatorDstAddress)
		if err != nil {
			return nil, err
		}
		delegatorAddress, err := AddressCodec.StringToBytes(dvvTriplet.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		balances, err := am.keeper.CompleteRedelegation(
			ctx,
			delegatorAddress,
			valSrcAddr,
			valDstAddr,
		)
		if err != nil {
			continue
		}

		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteRedelegation,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvvTriplet.DelegatorAddress),
				sdk.NewAttribute(types.AttributeKeySrcValidator, dvvTriplet.ValidatorSrcAddress),
				sdk.NewAttribute(types.AttributeKeyDstValidator, dvvTriplet.ValidatorDstAddress),
			),
		)
	}

	return validatorUpdates, nil
}

// Use custom UnbondAllMatureValidators
func (am AppModule) BlockValidatorUpdatesPatch2(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	validatorAddressCodec := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())
	AddressCodec := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	// Calculate validator set changes.
	//
	// NOTE: ApplyAndReturnValidatorSetUpdates has to come before
	// UnbondAllMatureValidatorQueue.
	// This fixes a bug when the unbonding period is instant (is the case in
	// some of the tests). The test expected the validator to be completely
	// unbonded after the Endblocker (go from Bonded -> Unbonding during
	// ApplyAndReturnValidatorSetUpdates and then Unbonding -> Unbonded during
	// UnbondAllMatureValidatorQueue).
	validatorUpdates, err := am.keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	if err != nil {
		return nil, err
	}

	// unbond all mature validators from the unbonding queue
	err = am.UnbondAllMatureValidatorsPatch2(ctx)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds, err := am.keeper.DequeueAllMatureUBDQueue(ctx, sdkCtx.BlockHeader().Time)
	if err != nil {
		return nil, err
	}

	for _, dvPair := range matureUnbonds {
		addr, err := validatorAddressCodec.StringToBytes(dvPair.ValidatorAddress)
		if err != nil {
			return nil, err
		}
		delegatorAddress, err := AddressCodec.StringToBytes(dvPair.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		balances, err := am.keeper.CompleteUnbonding(ctx, delegatorAddress, addr)
		if err != nil {
			continue
		}

		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbonding,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, dvPair.ValidatorAddress),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress),
			),
		)
	}

	// Remove all mature redelegations from the red queue.
	matureRedelegations, err := am.keeper.DequeueAllMatureRedelegationQueue(ctx, sdkCtx.BlockHeader().Time)
	if err != nil {
		return nil, err
	}

	for _, dvvTriplet := range matureRedelegations {
		valSrcAddr, err := validatorAddressCodec.StringToBytes(dvvTriplet.ValidatorSrcAddress)
		if err != nil {
			return nil, err
		}
		valDstAddr, err := validatorAddressCodec.StringToBytes(dvvTriplet.ValidatorDstAddress)
		if err != nil {
			return nil, err
		}
		delegatorAddress, err := AddressCodec.StringToBytes(dvvTriplet.DelegatorAddress)
		if err != nil {
			return nil, err
		}

		balances, err := am.keeper.CompleteRedelegation(
			ctx,
			delegatorAddress,
			valSrcAddr,
			valDstAddr,
		)
		if err != nil {
			continue
		}

		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteRedelegation,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvvTriplet.DelegatorAddress),
				sdk.NewAttribute(types.AttributeKeySrcValidator, dvvTriplet.ValidatorSrcAddress),
				sdk.NewAttribute(types.AttributeKeyDstValidator, dvvTriplet.ValidatorDstAddress),
			),
		)
	}

	return validatorUpdates, nil
}

// if the validator is not unbonding, set it to unbonding
func (am AppModule) UnbondAllMatureValidatorsPatch2(ctx context.Context) error {
	validatorAddressCodec := authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	blockTime := sdkCtx.BlockTime()
	blockHeight := sdkCtx.BlockHeight()

	// unbondingValIterator will contains all validator addresses indexed under
	// the ValidatorQueueKey prefix. Note, the entire index key is composed as
	// ValidatorQueueKey | timeBzLen (8-byte big endian) | timeBz | heightBz (8-byte big endian),
	// so it may be possible that certain validator addresses that are iterated
	// over are not ready to unbond, so an explicit check is required.
	unbondingValIterator, err := am.keeper.ValidatorQueueIterator(ctx, blockTime, blockHeight)
	if err != nil {
		return err
	}
	defer unbondingValIterator.Close()

	for ; unbondingValIterator.Valid(); unbondingValIterator.Next() {
		key := unbondingValIterator.Key()
		keyTime, keyHeight, err := types.ParseValidatorQueueKey(key)
		if err != nil {
			return fmt.Errorf("failed to parse unbonding key: %w", err)
		}

		// All addresses for the given key have the same unbonding height and time.
		// We only unbond if the height and time are less than the current height
		// and time.
		if keyHeight <= blockHeight && (keyTime.Before(blockTime) || keyTime.Equal(blockTime)) {
			addrs := types.ValAddresses{}
			if err = am.cdc.Unmarshal(unbondingValIterator.Value(), &addrs); err != nil {
				return err
			}

			for _, valAddr := range addrs.Addresses {
				addr, err := validatorAddressCodec.StringToBytes(valAddr)
				if err != nil {
					return err
				}
				val, err := am.keeper.GetValidator(ctx, addr)
				if err != nil {
					return errorsmod.Wrap(err, "validator in the unbonding queue was not found")
				}

				// if the validator is not unbonding, set it to unbonding
				if !val.IsUnbonding() {
					val.Status = types.Unbonding
					if err = am.keeper.SetValidator(ctx, val); err != nil {
						return err
					}
				}

				if val.UnbondingOnHoldRefCount == 0 {
					for _, id := range val.UnbondingIds {
						if err = am.keeper.DeleteUnbondingIndex(ctx, id); err != nil {
							return err
						}
					}

					val, err = am.keeper.UnbondingToUnbonded(ctx, val)
					if err != nil {
						return err
					}

					if val.GetDelegatorShares().IsZero() {
						str, err := validatorAddressCodec.StringToBytes(val.GetOperator())
						if err != nil {
							return err
						}
						if err = am.keeper.RemoveValidator(ctx, str); err != nil {
							return err
						}
					} else {
						// remove unbonding ids
						val.UnbondingIds = []uint64{}
					}

					// remove validator from queue
					if err = am.keeper.DeleteValidatorQueue(ctx, val); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
