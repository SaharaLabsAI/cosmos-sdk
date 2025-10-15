package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

const (
	// ValidatorBonusConfigKey is the key for storing validator bonus configuration
	ValidatorBonusConfigKey = "validator_bonus_config"
)

// SetValidatorBonusConfig sets the validator bonus configuration
func (k Keeper) SetValidatorBonusConfig(ctx context.Context, config types.ValidatorBonusConfig) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&config)
	if err != nil {
		return err
	}
	return store.Set([]byte(ValidatorBonusConfigKey), bz)
}

// GetValidatorBonusConfig gets the validator bonus configuration
func (k Keeper) GetValidatorBonusConfig(ctx context.Context) (types.ValidatorBonusConfig, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get([]byte(ValidatorBonusConfigKey))
	if err != nil {
		return types.ValidatorBonusConfig{}, err
	}

	if bz == nil {
		return types.DefaultValidatorBonusConfig(), nil
	}

	var config types.ValidatorBonusConfig
	err = k.cdc.Unmarshal(bz, &config)
	return config, err
}

// GetGenesisState returns the current genesis state including bonus configuration
func (k Keeper) GetGenesisState(ctx context.Context) (*types.GenesisState, error) {
	feePool, err := k.FeePool.Get(ctx)
	if err != nil {
		return nil, err
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Get all other genesis data (simplified for now)
	dwi := make([]types.DelegatorWithdrawInfo, 0)
	k.IterateDelegatorWithdrawAddrs(ctx, func(del, addr sdk.AccAddress) (stop bool) {
		dwi = append(dwi, types.DelegatorWithdrawInfo{
			DelegatorAddress: del.String(),
			WithdrawAddress:  addr.String(),
		})
		return false
	})

	pp, err := k.GetPreviousProposerConsAddr(ctx)
	if err != nil {
		return nil, err
	}

	outstanding := make([]types.ValidatorOutstandingRewardsRecord, 0)
	k.IterateValidatorOutstandingRewards(ctx,
		func(addr sdk.ValAddress, rewards types.ValidatorOutstandingRewards) (stop bool) {
			outstanding = append(outstanding, types.ValidatorOutstandingRewardsRecord{
				ValidatorAddress:   addr.String(),
				OutstandingRewards: rewards.Rewards,
			})
			return false
		},
	)

	acc := make([]types.ValidatorAccumulatedCommissionRecord, 0)
	k.IterateValidatorAccumulatedCommissions(ctx,
		func(addr sdk.ValAddress, commission types.ValidatorAccumulatedCommission) (stop bool) {
			acc = append(acc, types.ValidatorAccumulatedCommissionRecord{
				ValidatorAddress: addr.String(),
				Accumulated:      commission,
			})
			return false
		},
	)

	his := make([]types.ValidatorHistoricalRewardsRecord, 0)
	k.IterateValidatorHistoricalRewards(ctx,
		func(val sdk.ValAddress, period uint64, rewards types.ValidatorHistoricalRewards) (stop bool) {
			his = append(his, types.ValidatorHistoricalRewardsRecord{
				ValidatorAddress: val.String(),
				Period:           period,
				Rewards:          rewards,
			})
			return false
		},
	)

	cur := make([]types.ValidatorCurrentRewardsRecord, 0)
	k.IterateValidatorCurrentRewards(ctx,
		func(val sdk.ValAddress, rewards types.ValidatorCurrentRewards) (stop bool) {
			cur = append(cur, types.ValidatorCurrentRewardsRecord{
				ValidatorAddress: val.String(),
				Rewards:          rewards,
			})
			return false
		},
	)

	dels := make([]types.DelegatorStartingInfoRecord, 0)
	k.IterateDelegatorStartingInfos(ctx,
		func(val sdk.ValAddress, del sdk.AccAddress, info types.DelegatorStartingInfo) (stop bool) {
			dels = append(dels, types.DelegatorStartingInfoRecord{
				ValidatorAddress: val.String(),
				DelegatorAddress: del.String(),
				StartingInfo:     info,
			})
			return false
		},
	)

	slashes := make([]types.ValidatorSlashEventRecord, 0)
	k.IterateValidatorSlashEvents(ctx,
		func(val sdk.ValAddress, height uint64, event types.ValidatorSlashEvent) (stop bool) {
			slashes = append(slashes, types.ValidatorSlashEventRecord{
				ValidatorAddress:    val.String(),
				Height:              height,
				Period:              event.ValidatorPeriod,
				ValidatorSlashEvent: event,
			})
			return false
		},
	)

	// Get validator bonus configuration
	bonusConfig, err := k.GetValidatorBonusConfig(ctx)
	if err != nil {
		return nil, err
	}

	return types.NewGenesisState(params, feePool, dwi, pp, outstanding, acc, his, cur, dels, slashes, bonusConfig), nil
}
