package rest

import (
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
)

// Context defines the read-only request execution context supplied to
// handlers and middlewares. It provides getter and binding access to the
// incoming HTTP request without exposing internal mutating operations.
type Context = contracts.Context
