package config

import (
	"math"
	"testing"
)

const (
	testSampleRateBelowZero = -0.01
	testSampleRateAboveOne  = 1.0001
	testSampleRateNegative  = -0.5
	testSampleRateHighOut   = 1.5
)

// TestConfig_Validate_SampleRateValid tests valid boundary and interior sample rates.
func TestConfig_Validate_SampleRateValid(t *testing.T) {
	zeroRateCfg := validTestConfig()
	zeroRateCfg.SampleRate = MinSampleRate

	halfRateCfg := validTestConfig()
	halfRateCfg.SampleRate = testSampleRateHalf

	oneRateCfg := validTestConfig()
	oneRateCfg.SampleRate = MaxSampleRate

	runValidateCases(t, []configValidateCase{
		{cfg: zeroRateCfg, name: "minimum rate 0.0", expectedErr: ""},
		{cfg: halfRateCfg, name: "intermediate rate 0.5", expectedErr: ""},
		{cfg: oneRateCfg, name: "maximum rate 1.0", expectedErr: ""},
	})
}

// TestConfig_Validate_SampleRateInvalid tests out of bounds, NaN, and Inf rates.
func TestConfig_Validate_SampleRateInvalid(t *testing.T) {
	negCfg := validTestConfig()
	negCfg.SampleRate = testSampleRateBelowZero

	highCfg := validTestConfig()
	highCfg.SampleRate = testSampleRateAboveOne

	nanCfg := validTestConfig()
	nanCfg.SampleRate = math.NaN()

	infCfg := validTestConfig()
	infCfg.SampleRate = math.Inf(1)

	runValidateCases(t, []configValidateCase{
		{cfg: negCfg, name: "negative rate", expectedErr: "sample rate must be between 0.0 and 1.0"},
		{cfg: highCfg, name: "rate above 1.0", expectedErr: "sample rate must be between 0.0 and 1.0"},
		{cfg: nanCfg, name: "NaN rate", expectedErr: "sample rate must be between 0.0 and 1.0"},
		{cfg: infCfg, name: "+Inf rate", expectedErr: "sample rate must be between 0.0 and 1.0"},
	})
}
