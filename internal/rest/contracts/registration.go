package contracts

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
	Apis    []*ExportableAPI
	Version APIVersioning
}
