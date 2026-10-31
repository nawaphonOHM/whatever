package provider

import "testing"

func TestInitTracerProvider_JSONLogging_Disabled(t *testing.T) {
	assertDisabledProviderLogs(t, jsonProviderLogMode())
}
