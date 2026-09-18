package server

// validatedRoute is a fully resolved route ready for engine mount.
type validatedRoute struct {
	Handler     Handler
	Path        string
	Method      string
	Middlewares []Middleware
}
