package keeper

import (
	"context"
	"strconv"

	"cosmossdk.io/math"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerTokenomicsTypes "github.com/elys-network/elys/indexer/txs/tokenomics"
	indexerTypes "github.com/elys-network/elys/indexer/types"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/tokenomics/types"
)

func (k msgServer) CreateAirdrop(goCtx context.Context, msg *types.MsgCreateAirdrop) (*types.MsgCreateAirdropResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, found := k.GetAirdrop(ctx, msg.Intent)
	if found {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	airdrop := types.Airdrop{
		Authority: msg.Authority,
		Intent:    msg.Intent,
		Amount:    msg.Amount,
		Expiry:    msg.Expiry,
	}

	k.SetAirdrop(ctx, airdrop)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerTokenomicsTypes.MsgCreateAirdrop{
		Authority: msg.Authority,
		Intent:    msg.Intent,
		Amount:    msg.Amount,
		Expiry:    msg.Expiry,
	}, []string{msg.Authority})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgCreateAirdropResponse{}, nil
}

func (k msgServer) UpdateAirdrop(goCtx context.Context, msg *types.MsgUpdateAirdrop) (*types.MsgUpdateAirdropResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, found := k.GetAirdrop(ctx, msg.Intent)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg authority is the same as the current owner
	if msg.Authority != valFound.Authority {
		return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	airdrop := types.Airdrop{
		Authority: msg.Authority,
		Intent:    msg.Intent,
		Amount:    msg.Amount,
		Expiry:    msg.Expiry,
	}

	k.SetAirdrop(ctx, airdrop)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerTokenomicsTypes.MsgUpdateAirdrop{
		Authority: msg.Authority,
		Intent:    msg.Intent,
		Amount:    msg.Amount,
		Expiry:    msg.Expiry,
	}, []string{msg.Authority})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgUpdateAirdropResponse{}, nil
}

func (k msgServer) DeleteAirdrop(goCtx context.Context, msg *types.MsgDeleteAirdrop) (*types.MsgDeleteAirdropResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, found := k.GetAirdrop(ctx, msg.Intent)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg authority is the same as the current owner
	if msg.Authority != valFound.Authority {
		return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveAirdrop(ctx, msg.Intent)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerTokenomicsTypes.MsgDeleteAirdrop{
		Authority: msg.Authority,
		Intent:    msg.Intent,
	}, []string{msg.Authority})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgDeleteAirdropResponse{}, nil
}

func (k msgServer) ClaimAirdrop(goCtx context.Context, msg *types.MsgClaimAirdrop) (*types.MsgClaimAirdropResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	airdrop, found := k.GetAirdrop(ctx, msg.Sender)
	if !found {
		return nil, errors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg authority is the same as the current owner
	if msg.Sender != airdrop.Authority {
		return nil, errors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if ctx.BlockTime().Unix() > int64(airdrop.Expiry) {
		return nil, types.ErrAirdropExpired
	}

	// Add commitments
	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	commitments := k.commitmentKeeper.GetCommitments(ctx, sender)
	commitments.AddClaimed(sdk.NewCoin(ptypes.Eden, math.NewInt(int64(airdrop.Amount))))
	k.commitmentKeeper.SetCommitments(ctx, commitments)

	k.RemoveAirdrop(ctx, msg.Sender)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	// Convert committed tokens to indexer format
	indexerCommittedTokens := make([]*indexerTokenomicsTypes.CommittedTokens, len(commitments.CommittedTokens))
	for i, ct := range commitments.CommittedTokens {
		lockups := make([]indexerTokenomicsTypes.Lockup, len(ct.Lockups))
		for j, l := range ct.Lockups {
			lockups[j] = indexerTokenomicsTypes.Lockup{
				Amount:          l.Amount.String(),
				UnlockTimestamp: l.UnlockTimestamp,
			}
		}
		indexerCommittedTokens[i] = &indexerTokenomicsTypes.CommittedTokens{
			Denom:   ct.Denom,
			Amount:  ct.Amount.String(),
			Lockups: lockups,
		}
	}

	// Convert vesting tokens to indexer format
	indexerVestingTokens := make([]*indexerTokenomicsTypes.VestingTokens, len(commitments.VestingTokens))
	for i, vt := range commitments.VestingTokens {
		indexerVestingTokens[i] = &indexerTokenomicsTypes.VestingTokens{
			Denom:                vt.Denom,
			TotalAmount:          vt.TotalAmount.String(),
			ClaimedAmount:        vt.ClaimedAmount.String(),
			NumBlocks:            vt.NumBlocks,
			StartBlock:           vt.StartBlock,
			VestStartedTimestamp: vt.VestStartedTimestamp,
		}
	}

	indexer.QueueTransaction(ctx, indexerTokenomicsTypes.MsgClaimAirdrop{
		Sender: msg.Sender,
		AmountClaimed: indexerTypes.Token{
			Amount: strconv.FormatUint(airdrop.Amount, 10),
			Denom:  ptypes.Eden,
		},
		CommittedTokens: indexerCommittedTokens,
		CommitsClaimed: indexerTypes.Token{
			Amount: commitments.Claimed.AmountOf(ptypes.Eden).String(),
			Denom:  ptypes.Eden,
		},
		VestingTokens: indexerVestingTokens,
	}, []string{msg.Sender})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	return &types.MsgClaimAirdropResponse{}, nil
}
