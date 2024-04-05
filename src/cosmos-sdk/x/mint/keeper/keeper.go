package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/x/mint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the mint store
type Keeper struct {
	cdc              codec.BinaryCodec
	environment      appmodule.Environment
	stakingKeeper    types.StakingKeeper
	bankKeeper       types.BankKeeper
	logger           log.Logger
	feeCollectorName string
	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string

	Schema                  collections.Schema
	Params                  collections.Item[types.Params]
	Minter                  collections.Item[types.Minter]
	LastStakeCoinMintTime   collections.Item[math.Int]
	LastStakeCoinMintAmount collections.Item[math.Int]
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	env appmodule.Environment,
	sk types.StakingKeeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	feeCollectorName string,
	authority string,
) Keeper {
	// ensure mint module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("the x/%s module account has not been set", types.ModuleName))
	}

	sb := collections.NewSchemaBuilder(env.KVStoreService)
	k := Keeper{
		cdc:                     cdc,
		environment:             env,
		stakingKeeper:           sk,
		bankKeeper:              bk,
		logger:                  env.Logger,
		feeCollectorName:        feeCollectorName,
		authority:               authority,
		Params:                  collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Minter:                  collections.NewItem(sb, types.MinterKey, "minter", codec.CollValue[types.Minter](cdc)),
		LastStakeCoinMintAmount: collections.NewItem(sb, types.LastStakeCoinMintAmountKey, "last_stake_coin_mint_amount", sdk.IntValue),
		LastStakeCoinMintTime:   collections.NewItem(sb, types.LastStakeCoinMintTimeKey, "last_stake_coin_mint_time", sdk.IntValue),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

// GetAuthority returns the x/mint module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	return k.environment.Logger.With("module", "x/"+types.ModuleName)
}

// StakingTokenSupply implements an alias call to the underlying staking keeper's
// StakingTokenSupply to be used in BeginBlocker.
func (k Keeper) StakingTokenSupply(ctx context.Context) (math.Int, error) {
	return k.stakingKeeper.StakingTokenSupply(ctx)
}

// BondedRatio implements an alias call to the underlying staking keeper's
// BondedRatio to be used in BeginBlocker.
func (k Keeper) BondedRatio(ctx context.Context) (math.LegacyDec, error) {
	return k.stakingKeeper.BondedRatio(ctx)
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx context.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}

	return k.bankKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// AddCollectedFees implements an alias call to the underlying supply keeper's
// AddCollectedFees to be used in BeginBlocker.
func (k Keeper) AddCollectedFees(ctx context.Context, fees sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)
}

func (k Keeper) GetLastStakeCoinMintAmount(ctx context.Context) (math.Int, error) {
	fmt.Println("GetLastStakeCoinMintAmount")
	amount, err := k.LastStakeCoinMintAmount.Get(ctx)
	fmt.Println("Get amount: ", amount)
	if err != nil {
		return amount, err
	}
	return amount, nil
}

func (k Keeper) SetLastStakeCoinMintAmount(ctx context.Context, amount math.Int) error {
	fmt.Println("SetLastStakeCoinMintAmount")
	fmt.Println("Set amount: ", amount)
	return k.LastStakeCoinMintAmount.Set(ctx, amount)
}

func (k Keeper) GetLastStakeCoinMintTime(ctx context.Context) (math.Int, error) {
	fmt.Println("GetLastStakeCoinMintTime")
	t, err := k.LastStakeCoinMintTime.Get(ctx)
	fmt.Println("Get t: ", t)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (k Keeper) SetLastStakeCoinMintTime(ctx context.Context, t math.Int) error {
	fmt.Println("SetLastStakeCoinMintTime")
	fmt.Println("Set t: ", t)
	return k.LastStakeCoinMintTime.Set(ctx, t)
}

func (k Keeper) HasLastStakeCoinMintAmount(ctx context.Context) bool {
	ifHas, err := k.LastStakeCoinMintAmount.Has(ctx)
	if err != nil {
		return false
	}
	return ifHas
}

func (k Keeper) HasLastStakeCoinMintTime(ctx context.Context) bool {
	ifHas, err := k.LastStakeCoinMintTime.Has(ctx)
	if err != nil {
		return false
	}
	return ifHas
}
