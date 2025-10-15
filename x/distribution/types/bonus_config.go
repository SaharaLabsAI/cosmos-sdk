package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewValidatorBonusConfig creates a new ValidatorBonusConfig
func NewValidatorBonusConfig(whitelist []string, heightTo uint64, bonusPercent math.LegacyDec) ValidatorBonusConfig {
	return ValidatorBonusConfig{
		ValidatorWhitelist: whitelist,
		BonusHeightTo:      heightTo,
		BonusPercent:       bonusPercent,
	}
}

// DefaultValidatorBonusConfig returns default bonus configuration
func DefaultValidatorBonusConfig() ValidatorBonusConfig {
	return ValidatorBonusConfig{
		ValidatorWhitelist: []string{},
		BonusHeightTo:      0,
		BonusPercent:       math.LegacyZeroDec(),
	}
}

// IsEligibleForBonus checks if a validator is eligible for bonus at given height
func (cfg ValidatorBonusConfig) IsEligibleForBonus(validatorAddr string, currentHeight uint64) bool {
	// Check if current height is within bonus period
	if currentHeight > cfg.BonusHeightTo {
		return false
	}

	// Check if validator is in whitelist
	for _, addr := range cfg.ValidatorWhitelist {
		if addr == validatorAddr {
			return true
		}
	}
	return false
}

// CalculateBonus calculates the bonus amount for given reward
func (cfg ValidatorBonusConfig) CalculateBonus(baseReward sdk.DecCoins) sdk.DecCoins {
	if cfg.BonusPercent.IsZero() {
		return sdk.NewDecCoins()
	}

	bonusCoins := make(sdk.DecCoins, len(baseReward))
	for i, coin := range baseReward {
		bonusAmount := coin.Amount.Mul(cfg.BonusPercent).Quo(math.LegacyNewDec(100))
		bonusCoins[i] = sdk.NewDecCoinFromDec(coin.Denom, bonusAmount)
	}
	return bonusCoins
}
