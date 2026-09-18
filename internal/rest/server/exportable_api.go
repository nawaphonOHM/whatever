package server

// ExportableAPI defines a single API route endpoint.
type ExportableAPI struct {
	Path       Pathz
	Handler    Handler
	Middleware []Middleware
	Method     HTTPMethod
}

// RestAPIRegistration groups APIs under a common prefix and version.
type RestAPIRegistration struct {
	Prefix  Pathz
	Apis    []*ExportableAPI
	Version APIVersioning
}
