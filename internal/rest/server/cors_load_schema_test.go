package server

import "testing"

func TestLoadAndValidateCORS_MethodAndHeaderViolations(t *testing.T) {
	cases := []schemaTestCase{
		{
			name:    "unknown extra property",
			content: `{"allow-origins": ["https://example.com"], "unknown": "val"}`,
			wantErr: "Additional property unknown is not allowed",
		},
		{
			name:    "invalid HTTP method enum",
			content: `{"allow-methods": ["INVALID_METHOD"]}`,
			wantErr: "allow-methods.0 must be one of the following",
		},
		{
			name:    "PATCH method rejected",
			content: `{"allow-methods": ["PATCH"]}`,
			wantErr: "allow-methods.0 must be one of the following",
		},
		{
			name:    "invalid HTTP header enum",
			content: `{"allow-headers": ["Invalid-Header-XYZ"]}`,
			wantErr: "must be one of the following",
		},
		{
			name:    "lowercase header rejected",
			content: `{"allow-headers": ["content-type"]}`,
			wantErr: "must be one of the following",
		},
		{
			name:    "invalid request-headers enum",
			content: `{"request-headers": ["Bogus-Header-123"]}`,
			wantErr: "must be one of the following",
		},
		{
			name:    "invalid request-methods enum",
			content: `{"request-methods": ["UNKNOWN_METHOD"]}`,
			wantErr: "request-methods.0 must be one of the following",
		},
		{
			name:    "mixed valid and invalid allow-methods at index 1",
			content: `{"allow-methods": ["GET", "PATCH"]}`,
			wantErr: "allow-methods.1 must be one of the following",
		},
		{
			name:    "mixed valid and invalid request-methods at index 1",
			content: `{"request-methods": ["POST", "INVALID_METHOD"]}`,
			wantErr: "request-methods.1 must be one of the following",
		},
		{
			name:    "lowercase allow-methods rejected",
			content: `{"allow-methods": ["get"]}`,
			wantErr: "allow-methods.0 must be one of the following",
		},
		{
			name:    "lowercase request-methods rejected",
			content: `{"request-methods": ["post"]}`,
			wantErr: "request-methods.0 must be one of the following",
		},
	}
	runSchemaTestCases(t, cases)
}
