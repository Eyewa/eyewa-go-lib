package exporter

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

type stdOutExporter struct {
	exporter *stdouttrace.Exporter
}

type otelCollectorExporter struct {
	exporter *otlptrace.Exporter
}

// Exporter is a backwards compatible interface for an otlptrace.Exporter.
type Exporter interface {
	// Start starts the exporter
	Start(ctx context.Context) error

	// ExportSpans exports spans to a span destination.
	ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error

	// Shutdown shuts down the connection to a span destination.
	Shutdown(ctx context.Context) error
}
