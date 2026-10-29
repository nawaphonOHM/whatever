package rest

import (
	"fmt"
	"testing"
)

const testHostURLFormat = "http://%s:%d"

func testPublicEnvIgnoredRawBlueprint(t *testing.T) {
	port := getFreePort(t)
	setStartEnv(t, port)
	t.Setenv(envAPIGroup, testGroupStd)

	bp := makeGroupBlueprint(APIGroupRaw, testCustomPath)
	errCh := startRESTAsync(bp)
	baseURL := fmt.Sprintf(testHostURLFormat, testHost, port)

	assertStatusOK(t, baseURL+testCustomPath)
	signalAndWait(t, errCh)
}

func testPublicEnvIgnoredStdBlueprint(t *testing.T) {
	port := getFreePort(t)
	setStartEnv(t, port)
	t.Setenv(envAPIGroup, testGroupRaw)

	bp := makeGroupBlueprint(APIGroupStandard, testCustomPath)
	errCh := startRESTAsync(bp)
	baseURL := fmt.Sprintf(testHostURLFormat, testHost, port)

	assertStatusOK(t, baseURL+"/api"+testCustomPath)
	signalAndWait(t, errCh)
}

func testPublicEnvIgnoredDefaultBlueprint(t *testing.T) {
	port := getFreePort(t)
	setStartEnv(t, port)
	t.Setenv(envAPIGroup, testGroupRaw)

	bp := makeGroupBlueprint("", testCustomPath)
	errCh := startRESTAsync(bp)
	baseURL := fmt.Sprintf(testHostURLFormat, testHost, port)

	assertStatusOK(t, baseURL+"/api"+testCustomPath)
	signalAndWait(t, errCh)
}

func testPublicEnvIgnoredInvalidEnv(t *testing.T) {
	port := getFreePort(t)
	setStartEnv(t, port)
	t.Setenv(envAPIGroup, "invalid_environment_group")

	bp := makeGroupBlueprint(APIGroupRaw, testCustomPath)
	errCh := startRESTAsync(bp)
	baseURL := fmt.Sprintf(testHostURLFormat, testHost, port)

	assertStatusOK(t, baseURL+testCustomPath)
	signalAndWait(t, errCh)
}

func TestStartREST_APIGroupEnvIgnored(t *testing.T) {
	t.Run("StdEnv_RawBlueprint_Ignored", testPublicEnvIgnoredRawBlueprint)
	t.Run("RawEnv_StdBlueprint_Ignored", testPublicEnvIgnoredStdBlueprint)
	t.Run("RawEnv_DefaultBlueprint_Ignored", testPublicEnvIgnoredDefaultBlueprint)
	t.Run("InvalidEnv_Blueprint_Ignored", testPublicEnvIgnoredInvalidEnv)
}
