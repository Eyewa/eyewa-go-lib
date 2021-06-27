package exporter

import (
	"context"

	"go.opentelemetry.io/otel/sdk/trace"
)

type stdOutExporter struct {
	exporter stdOutExporterIface
}

type otelCollectorExporter struct {
	exporter otelExporterIface
	endpoint string
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

type stdOutExporterIface interface {
	// ExportSpans exports spans to a span destination.
	ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error

	// Shutdown shuts down the connection to a span destination.
	Shutdown(ctx context.Context) error
}

type otelExporterIface interface {
	// Start starts the exporter
	Start(ctx context.Context) error

	// ExportSpans exports spans to a span destination.
	ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error

	// Shutdown shuts down the connection to a span destination.
	Shutdown(ctx context.Context) error
}
