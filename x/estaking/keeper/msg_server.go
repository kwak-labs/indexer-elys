package keeper

import (
	"context"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerEstakingTypes "github.com/elys-network/elys/indexer/txs/estaking"
	indexerTypes "github.com/elys-network/elys/indexer/types"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/elys-network/elys/x/estaking/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	k.SetParams(ctx, req.Params)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerEstakingTypes.MsgUpdateParams{
		Authority: req.Authority,
		Params:    req.Params,
	}, []string{})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdateParamsResponse{}, nil
}

func (k msgServer) WithdrawReward(goCtx context.Context, msg *types.MsgWithdrawReward) (*types.MsgWithdrawRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	delAddr := sdk.MustAccAddressFromBech32(msg.DelegatorAddress)
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	amount, err := k.distrKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
	if err != nil {
		return nil, err
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	tokens := make([]indexerTypes.Token, len(amount))
	for i, coin := range amount {
		tokens[i] = indexerTypes.Token{
			Amount: coin.Amount.String(),
			Denom:  coin.Denom,
		}
	}

	indexer.QueueTransaction(ctx, indexerEstakingTypes.MsgWithdrawReward{
		DelegatorAddress: msg.DelegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
		Amount:           tokens,
	}, []string{})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtWithdrawReward,
			sdk.NewAttribute(types.AttributeDelegatorAddress, msg.DelegatorAddress),
			sdk.NewAttribute(types.AttributeValidatorAddress, msg.ValidatorAddress),
			sdk.NewAttribute(types.AttributeAmount, amount.String()),
		),
	})
	return &types.MsgWithdrawRewardResponse{Amount: amount}, nil
}

func (k msgServer) WithdrawElysStakingRewards(goCtx context.Context, msg *types.MsgWithdrawElysStakingRewards) (*types.MsgWithdrawElysStakingRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	delAddr := sdk.MustAccAddressFromBech32(msg.DelegatorAddress)

	var amount sdk.Coins
	var err error = nil
	var rewards = sdk.Coins{}
	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	var validatorAddresses []string
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */
	iterateError := k.Keeper.Keeper.IterateDelegations(ctx, delAddr, func(index int64, del stakingtypes.DelegationI) (stop bool) {
		valAddr, errB := sdk.ValAddressFromBech32(del.GetValidatorAddr())
		if errB != nil {
			err = errB
			return true
		}
		amount, err = k.distrKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
		if err != nil {
			return true
		}
		rewards = rewards.Add(amount...)
		/* *************************************************************************** */
		/* Start of kwak-indexer node implementation*/
		validatorAddresses = append(validatorAddresses, valAddr.String())
		/* End of kwak-indexer node implementation*/
		/* *************************************************************************** */

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.TypeEvtWithdrawReward,
				sdk.NewAttribute(types.AttributeDelegatorAddress, msg.DelegatorAddress),
				sdk.NewAttribute(types.AttributeValidatorAddress, valAddr.String()),
				sdk.NewAttribute(types.AttributeAmount, amount.String()),
			),
		})
		return false
	})
	if iterateError != nil {
		return nil, iterateError
	}
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	tokens := make([]indexerTypes.Token, len(rewards))
	for i, coin := range rewards {
		tokens[i] = indexerTypes.Token{
			Amount: coin.Amount.String(),
			Denom:  coin.Denom,
		}
	}

	indexer.QueueTransaction(ctx, indexerEstakingTypes.MsgWithdrawElysStakingRewards{
		DelegatorAddress: msg.DelegatorAddress,
		Validators:       validatorAddresses,
		Amount:           tokens,
	}, []string{msg.DelegatorAddress})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgWithdrawElysStakingRewardsResponse{Amount: rewards}, nil
}

func (k Keeper) WithdrawAllRewards(goCtx context.Context, msg *types.MsgWithdrawAllRewards) (*types.MsgWithdrawAllRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	delAddr := sdk.MustAccAddressFromBech32(msg.DelegatorAddress)
	var amount sdk.Coins
	var err error = nil
	var rewards = sdk.Coins{}
	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	var validatorAddresses []string
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	err = k.IterateDelegations(ctx, delAddr, func(index int64, del stakingtypes.DelegationI) (stop bool) {
		valAddr, errB := sdk.ValAddressFromBech32(del.GetValidatorAddr())
		if errB != nil {
			err = errB
			return true
		}
		amount, err = k.distrKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
		if err != nil {
			return true
		}
		rewards = rewards.Add(amount...)

		/* *************************************************************************** */
		/* Start of kwak-indexer node implementation*/
		validatorAddresses = append(validatorAddresses, valAddr.String())
		/* End of kwak-indexer node implementation*/
		/* *************************************************************************** */

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.TypeEvtWithdrawReward,
				sdk.NewAttribute(types.AttributeDelegatorAddress, msg.DelegatorAddress),
				sdk.NewAttribute(types.AttributeValidatorAddress, valAddr.String()),
				sdk.NewAttribute(types.AttributeAmount, amount.String()),
			),
		})
		return false
	})
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	tokens := make([]indexerTypes.Token, len(rewards))
	for i, coin := range rewards {
		tokens[i] = indexerTypes.Token{
			Amount: coin.Amount.String(),
			Denom:  coin.Denom,
		}
	}

	indexer.QueueTransaction(ctx, indexerEstakingTypes.MsgWithdrawAllRewards{
		DelegatorAddress: msg.DelegatorAddress,
		Validators:       validatorAddresses,
		Amount:           tokens,
	}, []string{msg.DelegatorAddress})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgWithdrawAllRewardsResponse{Amount: rewards}, nil
}
