package config

import (
	"testing"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace/noop"
)

const (
	testValTrue  = "true"
	testValFalse = "false"
	testValGCP   = "gcp"
	testValAWS   = "aws"

	testSampleRateLow   = 0.001
	testSampleRateRatio = 0.1
)

func setupActiveTracerWithSampler(t *testing.T, sampler sdktrace.Sampler) {
	t.Helper()
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
		sdktrace.WithSampler(sampler),
	)
	otel.SetTracerProvider(tp)
	t.Cleanup(func() {
		otel.SetTracerProvider(noop.NewTracerProvider())
	})
}

func setupActiveTracer(t *testing.T) {
	t.Helper()
	setupActiveTracerWithSampler(t, sdktrace.AlwaysSample())
}

func setupNoopTracer(t *testing.T) {
	t.Helper()
	otel.SetTracerProvider(noop.NewTracerProvider())
	t.Cleanup(func() {
		otel.SetTracerProvider(noop.NewTracerProvider())
	})
}
