package middleware

// Context key constants for trace and span correlation.
const (
	TraceIDKey     = "TraceID"
	SpanIDKey      = "SpanID"
	TraceIDAttrKey = "trace_id"
	SpanIDAttrKey  = "span_id"
)

// Header constants for W3C distributed trace context.
const (
	HeaderTraceParent = "traceparent"
	HeaderTraceState  = "tracestate"
)

// Standard OpenTelemetry semantic attribute keys for HTTP server spans.
const (
	AttrHTTPMethod     = "http.method"
	AttrHTTPRoute      = "http.route"
	AttrHTTPStatusCode = "http.status_code"
	AttrClientAddress  = "client.address"
	AttrHTTPTarget     = "http.target"
)

// DefaultTracerName is the default tracer instrumentation name.
const DefaultTracerName = "github.com/nawaphonOHM/whatever/internal/opentelemetry/middleware"

// statusServerError is the minimum status code representing a server error.
const statusServerError = 500
