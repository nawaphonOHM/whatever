package config

import (
	"testing"
	"time"
)

const (
	testPortZero        = 0
	testPortNegative    = -1
	testPortMin         = 1
	testPortStandard    = 27017
	testPortMax         = 65535
	testPortOverflow    = 65536
	testPortLarge       = 70000
	testLimitMaxPool    = 10
	testLimitMinPool    = 20
	errPortRangeMessage = "mongodb port must be between 1 and 65535"
)

func buildInvalidPortCases() []*configValidateCase {
	c0 := validTestConfig()
	c0.Port = testPortZero
	cNeg := validTestConfig()
	cNeg.Port = testPortNegative
	cOver := validTestConfig()
	cOver.Port = testPortOverflow
	cLarge := validTestConfig()
	cLarge.Port = testPortLarge

	return []*configValidateCase{
		{cfg: c0, name: "port zero on standard scheme", expectedErr: errPortRangeMessage},
		{cfg: cNeg, name: "negative port", expectedErr: errPortRangeMessage},
		{cfg: cOver, name: "port overflow 65536", expectedErr: errPortRangeMessage},
		{cfg: cLarge, name: "port 70000", expectedErr: errPortRangeMessage},
	}
}

func buildValidPortCases() []*configValidateCase {
	cMin := validTestConfig()
	cMin.Port = testPortMin
	cMax := validTestConfig()
	cMax.Port = testPortMax
	cStd := validTestConfig()
	cStd.Port = testPortStandard

	cSrvZero := validTestConfig()
	cSrvZero.Protocol = ProtocolMongoDBSrv
	cSrvZero.Port = 0

	return []*configValidateCase{
		{cfg: cMin, name: "valid min boundary port 1", expectedErr: ""},
		{cfg: cMax, name: "valid max boundary port 65535", expectedErr: ""},
		{cfg: cStd, name: "valid standard port 27017", expectedErr: ""},
		{cfg: cSrvZero, name: "valid port 0 on mongodb+srv", expectedErr: ""},
	}
}

func buildInvalidTimeoutCases() []*configValidateCase {
	cNegConn := validTestConfig()
	cNegConn.ConnectTimeout = -1 * time.Second
	cNegServer := validTestConfig()
	cNegServer.ServerSelectionTimeout = -1 * time.Second
	cNegSocket := validTestConfig()
	cNegSocket.SocketTimeout = -1 * time.Second
	cNegIdle := validTestConfig()
	cNegIdle.MaxConnIdleTime = -1 * time.Second

	return []*configValidateCase{
		{cfg: cNegConn, name: "negative connect timeout", expectedErr: "mongodb connect timeout cannot be negative"},
		{
			cfg:         cNegServer,
			name:        "negative server selection timeout",
			expectedErr: "mongodb server selection timeout cannot be negative",
		},
		{cfg: cNegSocket, name: "negative socket timeout", expectedErr: "mongodb socket timeout cannot be negative"},
		{
			cfg:         cNegIdle,
			name:        "negative max conn idle time",
			expectedErr: "mongodb max conn idle time cannot be negative",
		},
	}
}

func buildInvalidPoolCases() []*configValidateCase {
	cZeroMaxPool := validTestConfig()
	cZeroMaxPool.MaxPoolSize = 0

	cMinExceedsMax := validTestConfig()
	cMinExceedsMax.MaxPoolSize = testLimitMaxPool
	cMinExceedsMax.MinPoolSize = testLimitMinPool

	return []*configValidateCase{
		{
			cfg:         cZeroMaxPool,
			name:        "max pool size zero",
			expectedErr: "mongodb max pool size must be greater than 0",
		},
		{
			cfg:         cMinExceedsMax,
			name:        "min pool size exceeds max pool size",
			expectedErr: "mongodb min pool size cannot exceed max pool size",
		},
	}
}

// TestConfig_Validate_Limits tests port boundary, timeout, and pool limit validation rules.
func TestConfig_Validate_Limits(t *testing.T) {
	runValidateCases(t, buildInvalidPortCases())
	runValidateCases(t, buildValidPortCases())
	runValidateCases(t, buildInvalidTimeoutCases())
	runValidateCases(t, buildInvalidPoolCases())
}
