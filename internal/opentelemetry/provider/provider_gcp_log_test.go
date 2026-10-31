package provider

import "testing"

func TestInitTracerProvider_GCPLogging_Active(t *testing.T) {
	assertActiveProviderLogs(t, gcpProviderLogMode())
}

func TestInitTracerProvider_GCPLogging_Disabled(t *testing.T) {
	assertDisabledProviderLogs(t, gcpProviderLogMode())
}

func TestInitTracerProvider_GCPLogging_ValidationError(t *testing.T) {
	assertProviderValidationFailure(t, gcpProviderLogMode())
}
