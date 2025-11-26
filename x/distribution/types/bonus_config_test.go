package types_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/require"
)

func TestValidatorBonusConfig(t *testing.T) {
	// Test default configuration
	config := types.DefaultValidatorBonusConfig()
	require.Empty(t, config.ValidatorWhitelist)
	require.Equal(t, uint64(0), config.BonusHeightTo)
	require.True(t, config.BonusPercent.IsZero())

	// Test new configuration
	whitelist := []string{"cosmosvaloper1abc123", "cosmosvaloper1def456"}
	heightTo := uint64(1000)
	bonusPercent := math.LegacyNewDecWithPrec(10, 0) // 10%

	config = types.NewValidatorBonusConfig(whitelist, heightTo, bonusPercent)
	require.Equal(t, whitelist, config.ValidatorWhitelist)
	require.Equal(t, heightTo, config.BonusHeightTo)
	require.Equal(t, bonusPercent, config.BonusPercent)

	// Test eligibility check
	testCases := []struct {
		name           string
		validatorAddr  string
		currentHeight  uint64
		expectedResult bool
	}{
		{
			name:           "validator in whitelist, within height limit",
			validatorAddr:  "cosmosvaloper1abc123",
			currentHeight:  500,
			expectedResult: true,
		},
		{
			name:           "validator not in whitelist, within height limit",
			validatorAddr:  "cosmosvaloper1xyz789",
			currentHeight:  500,
			expectedResult: false,
		},
		{
			name:           "validator in whitelist, beyond height limit",
			validatorAddr:  "cosmosvaloper1abc123",
			currentHeight:  1500,
			expectedResult: false,
		},
		{
			name:           "validator not in whitelist, beyond height limit",
			validatorAddr:  "cosmosvaloper1xyz789",
			currentHeight:  1500,
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := config.IsEligibleForBonus(tc.validatorAddr, tc.currentHeight)
			require.Equal(t, tc.expectedResult, result)
		})
	}

	// Test bonus calculation
	baseReward := sdk.NewDecCoins(
		sdk.NewDecCoin("stake", math.NewInt(1000)),
		sdk.NewDecCoin("utoken", math.NewInt(500)),
	)

	bonus := config.CalculateBonus(baseReward)
	expectedBonus := sdk.NewDecCoins(
		sdk.NewDecCoin("stake", math.NewInt(100)), // 10% of 1000
		sdk.NewDecCoin("utoken", math.NewInt(50)), // 10% of 500
	)

	require.Equal(t, expectedBonus, bonus)

	// Test zero bonus percent
	configZero := types.NewValidatorBonusConfig(whitelist, heightTo, math.LegacyZeroDec())
	bonusZero := configZero.CalculateBonus(baseReward)
	require.True(t, bonusZero.IsZero())
}
