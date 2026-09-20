package contracts

// Handler defines route handlers returning a standardized Response.
type Handler func(Context) Response

// Middleware defines the function signature for route middlewares.
type Middleware func(Context)
