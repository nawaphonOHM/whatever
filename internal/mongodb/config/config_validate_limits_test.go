package config

import (
	"testing"
)

const (
	testPortZero        = 0
	testPortNegative    = -1
	testPortMin         = 1
	testPortStandard    = 27017
	testPortMax         = 65535
	testPortOverflow    = 65536
	testPortLarge       = 70000
	errPortRangeMessage = "mongodb port must be between 1 and 65535"
)

func buildInvalidPortCases() []configValidateCase {
	c0 := validTestConfig()
	c0.Port = testPortZero
	cNeg := validTestConfig()
	cNeg.Port = testPortNegative
	cOver := validTestConfig()
	cOver.Port = testPortOverflow
	cLarge := validTestConfig()
	cLarge.Port = testPortLarge

	return []configValidateCase{
		{cfg: c0, name: "port zero", expectedErr: errPortRangeMessage},
		{cfg: cNeg, name: "negative port", expectedErr: errPortRangeMessage},
		{cfg: cOver, name: "port overflow 65536", expectedErr: errPortRangeMessage},
		{cfg: cLarge, name: "port 70000", expectedErr: errPortRangeMessage},
	}
}

func buildValidPortCases() []configValidateCase {
	cMin := validTestConfig()
	cMin.Port = testPortMin
	cMax := validTestConfig()
	cMax.Port = testPortMax
	cStd := validTestConfig()
	cStd.Port = testPortStandard

	return []configValidateCase{
		{cfg: cMin, name: "valid min boundary port 1", expectedErr: ""},
		{cfg: cMax, name: "valid max boundary port 65535", expectedErr: ""},
		{cfg: cStd, name: "valid standard port 27017", expectedErr: ""},
	}
}

// TestConfig_Validate_Limits tests port boundary and limit validation rules.
func TestConfig_Validate_Limits(t *testing.T) {
	runValidateCases(t, buildInvalidPortCases())
	runValidateCases(t, buildValidPortCases())
}
