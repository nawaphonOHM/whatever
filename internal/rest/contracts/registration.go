package contracts

// APIGroup represents the API registration group mode.
type APIGroup string

const (
	// APIGroupStandard prefixes routes with /api and version segments (/api/v[N]).
	APIGroupStandard APIGroup = "standard"
	// APIGroupRaw registers routes directly with prefix and path without /api or version segments.
	APIGroupRaw APIGroup = "raw"
)

// ExportableAPI defines a single API route endpoint.
type ExportableAPI struct {
	Path       Pathz
	Handler    Handler
	Middleware []Middleware
	Method     HTTPMethod
}

// RRestAPIRegistration groups APIs under a common prefix and version.
type RRestAPIRegistration struct {
	Prefix  Pathz
	Group   APIGroup
	Apis    []*ExportableAPI
	Version APIVersioning
}
