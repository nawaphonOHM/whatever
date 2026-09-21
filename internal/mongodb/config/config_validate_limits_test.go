package config

import (
	"testing"
	"time"
)

const (
	testSmallPoolMax = uint64(5)
	testMinPoolLimit = uint64(10)
	testZeroPoolMax  = uint64(0)
	testNegDuration  = -1 * time.Second
)

func buildTimeoutCases() []configValidateCase {
	negConnect := validTestConfig()
	negConnect.ConnectTimeout = testNegDuration

	negServerSel := validTestConfig()
	negServerSel.ServerSelectionTimeout = testNegDuration

	negSocket := validTestConfig()
	negSocket.SocketTimeout = testNegDuration

	negIdle := validTestConfig()
	negIdle.MaxConnIdleTime = testNegDuration

	return []configValidateCase{
		{
			cfg:         negConnect,
			name:        "negative connect timeout",
			expectedErr: "connect timeout cannot be negative",
		},
		{
			cfg:         negServerSel,
			name:        "negative server selection timeout",
			expectedErr: "server selection timeout cannot be negative",
		},
		{
			cfg:         negSocket,
			name:        "negative socket timeout",
			expectedErr: "socket timeout cannot be negative",
		},
		{
			cfg:         negIdle,
			name:        "negative max conn idle time",
			expectedErr: "max conn idle time cannot be negative",
		},
	}
}

func buildPoolCases() []configValidateCase {
	invalidPool := validTestConfig()
	invalidPool.MaxPoolSize = testSmallPoolMax
	invalidPool.MinPoolSize = testMinPoolLimit

	zeroMaxPool := validTestConfig()
	zeroMaxPool.MaxPoolSize = testZeroPoolMax
	zeroMaxPool.MinPoolSize = testMinPoolLimit

	return []configValidateCase{
		{
			cfg:         invalidPool,
			name:        "min pool size greater than max pool size",
			expectedErr: "min pool size cannot be greater than max pool size",
		},
		{
			cfg:         zeroMaxPool,
			name:        "min pool size allowed when max pool size is zero",
			expectedErr: "",
		},
	}
}

// TestConfig_Validate_Limits tests timeout and connection pool validation rules.
func TestConfig_Validate_Limits(t *testing.T) {
	runValidateCases(t, buildTimeoutCases())
	runValidateCases(t, buildPoolCases())
}
