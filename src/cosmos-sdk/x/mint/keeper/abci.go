package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/core/event"
	"cosmossdk.io/math"
	"cosmossdk.io/x/mint/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker mints new tokens for the previous block.
func (k Keeper) BeginBlocker(ctx context.Context, ic types.InflationCalculationFn) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// fetch stored minter & params
	minter, err := k.Minter.Get(ctx)
	if err != nil {
		return err
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	// recalculate inflation rate
	totalStakingSupply, err := k.StakingTokenSupply(ctx)
	if err != nil {
		return err
	}

	bondedRatio, err := k.BondedRatio(ctx)
	if err != nil {
		return err
	}

	minter.Inflation = ic(ctx, minter, params, bondedRatio)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalStakingSupply)
	if err = k.Minter.Set(ctx, minter); err != nil {
		return err
	}

	// mint coins, update supply

	// edit START
	// @author: cxgd
	// @date: 2024-04-04
	// @description: This block of code is responsible for minting new StakeCoins.
	// The minting process is based on the time elapsed since the last minting operation.

	// Get the amount of coin to be minted for this block
	mintedCoin := minter.BlockProvision(params)
	var mintedCoins sdk.Coins

	// Get the current block's information
	headerInfo := k.environment.HeaderService.GetHeaderInfo(ctx)

	// Get the current block's time
	currentTime := headerInfo.Time
	fmt.Print("currentTime: ", currentTime)
	// Get the time of the last StakeCoin minting operation
	// Check if LastStakeCoinMintAmount has been set
	if !k.HasLastStakeCoinMintAmount(ctx) {
		k.SetLastStakeCoinMintAmount(ctx, math.NewInt(0))
	}

	// Check if LastStakeCoinMintTime has been set
	if !k.HasLastStakeCoinMintTime(ctx) {
		k.SetLastStakeCoinMintTime(ctx, math.NewInt(0))
	}
	lastStakeCoinMintTimeInt, err := k.GetLastStakeCoinMintTime(ctx)
	if err != nil {
		return err
	}
	lastStakeCoinMintTimeInt64 := lastStakeCoinMintTimeInt.Int64()
	lastStakeCoinMintTime := time.Unix(lastStakeCoinMintTimeInt64, 0)
	fmt.Print("lastStakeCoinMintTime: ", lastStakeCoinMintTime)
	// Calculate the duration since the last StakeCoin minting operation
	// If this is the first minting operation, set duration to a value greater than 7 days
	duration := 8 * 24 * time.Hour
	if !lastStakeCoinMintTime.IsZero() {
		// Otherwise, calculate the duration since the last minting operation
		duration = currentTime.Sub(lastStakeCoinMintTime)
	}
	fmt.Print("duration: ", duration)

	// Define the multiplier for minting coins
	coinMintMultiplier := math.NewInt(2)

	// If the duration is greater than or equal to 7 days, mint new StakeCoins
	// for test
	//if duration >= 30*time.Second {
	if duration >= 7*24*time.Hour {
		fmt.Print("Minting new StakeCoins")
		// Get the amount of StakeCoin minted in the last operation
		lastStakeCoinMintAmount, err := k.GetLastStakeCoinMintAmount(ctx)
		if err != nil {
			return err
		}
		// Determine the amount of StakeCoin to be minted
		stakeCoinMintAmount := math.NewInt(1)
		if !lastStakeCoinMintAmount.IsZero() {
			// If this is not the first minting operation, the amount of StakeCoin to be minted is twice the amount minted in the last operation
			stakeCoinMintAmount = lastStakeCoinMintAmount.Mul(coinMintMultiplier)
		}
		fmt.Print("stakeCoinMintAmount: ", stakeCoinMintAmount)
		// Mint new StakeCoins
		stakeCoin := sdk.NewCoin(sdk.StakingBondDenom, stakeCoinMintAmount)
		mintedCoins = sdk.NewCoins(stakeCoin)

		// Update the last mint amount and time
		k.SetLastStakeCoinMintAmount(ctx, stakeCoinMintAmount)
		k.SetLastStakeCoinMintTime(ctx, math.NewInt(currentTime.Unix()))
	} else {
		// If the duration is less than 7 days, no new StakeCoins are minted
		mintedCoins = sdk.NewCoins(mintedCoin)
	}
	// edit END

	// mintedCoin := minter.BlockProvision(params)
	// mintedCoins := sdk.NewCoins(mintedCoin)
	err = k.MintCoins(ctx, mintedCoins)
	if err != nil {
		return err
	}

	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, mintedCoins)
	if err != nil {
		return err
	}

	if mintedCoin.Amount.IsInt64() {
		defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
	}

	return k.environment.EventService.EventManager(ctx).EmitKV(
		types.EventTypeMint,
		event.NewAttribute(types.AttributeKeyBondedRatio, bondedRatio.String()),
		event.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
		event.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
		event.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
	)
}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
	// your logic here
}
