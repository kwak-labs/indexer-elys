package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer "github.com/elys-network/elys/indexer"
	indexerStableStakeTypes "github.com/elys-network/elys/indexer/txs/stablestake"

	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k msgServer) Bond(goCtx context.Context, msg *types.MsgBond) (*types.MsgBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	redemptionRate := k.GetRedemptionRate(ctx)

	depositDenom := k.GetDepositDenom(ctx)
	depositCoin := sdk.NewCoin(depositDenom, msg.Amount)
	err := k.bk.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, sdk.Coins{depositCoin})
	if err != nil {
		return nil, err
	}

	shareDenom := types.GetShareDenom()
	// Initial case
	if redemptionRate.IsZero() {
		redemptionRate = sdkmath.LegacyOneDec()
	}
	shareAmount := depositCoin.Amount.ToLegacyDec().Quo(redemptionRate).RoundInt()
	shareCoins := sdk.NewCoins(sdk.NewCoin(shareDenom, shareAmount))

	err = k.bk.MintCoins(ctx, types.ModuleName, shareCoins)
	if err != nil {
		return nil, err
	}

	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, shareCoins)
	if err != nil {
		return nil, err
	}

	_, found := k.assetProfileKeeper.GetEntry(ctx, shareDenom)
	if !found {
		// Set an entity to assetprofile
		entry := assetprofiletypes.Entry{
			Authority:                authtypes.NewModuleAddress(types.ModuleName).String(),
			BaseDenom:                shareDenom,
			Decimals:                 ptypes.BASE_DECIMAL,
			Denom:                    shareDenom,
			Path:                     "",
			IbcChannelId:             "",
			IbcCounterpartyChannelId: "",
			DisplayName:              shareDenom,
			DisplaySymbol:            "",
			Network:                  "",
			Address:                  "",
			ExternalSymbol:           "",
			TransferLimit:            "",
			Permissions:              make([]string, 0),
			UnitDenom:                "",
			IbcCounterpartyDenom:     "",
			IbcCounterpartyChainId:   "",
			CommitEnabled:            true,
			WithdrawEnabled:          true,
		}

		k.assetProfileKeeper.SetEntry(ctx, entry)
	}

	// Commit LP token
	err = k.commitmentKeeper.CommitLiquidTokens(ctx, creator, shareDenom, shareAmount, 0)
	if err != nil {
		return nil, err
	}

	params.TotalValue = params.TotalValue.Add(msg.Amount)
	k.SetParams(ctx, params)

	/* *************************************************************************** */
	/* Start of kwak-indexer node implementation*/
	indexer.QueueTransaction(ctx, indexerStableStakeTypes.MsgBond{
		Creator:        msg.Creator,
		Amount:         msg.Amount.String(),
		DepositDenom:   depositDenom,
		ShareAmount:    shareAmount.String(),
		ShareDenom:     shareDenom,
		RedemptionRate: redemptionRate.String(),
	}, []string{msg.Creator})
	/* End of kwak-indexer node implementation*/
	/* *************************************************************************** */

	if k.hooks != nil {
		err = k.hooks.AfterBond(ctx, creator, shareAmount)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgBondResponse{}, nil
}
