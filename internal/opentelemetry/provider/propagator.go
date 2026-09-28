package provider

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// NewCompositePropagator creates a W3C TraceContext and Baggage composite propagator.
func NewCompositePropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

// RegisterPropagators registers the default composite propagator globally.
func RegisterPropagators() {
	otel.SetTextMapPropagator(NewCompositePropagator())
}
