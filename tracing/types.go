package tracing

import (
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Config is the tracing environment configuration.
type Config struct {
	ServiceName             string `mapstructure:"service_name"`
	TracingExporterEndpoint string `mapstructure:"tracing_exporter_endpoint"`
	TracingSecureExporter   bool   `mapstructure:"tracing_secure_exporter"`
	TracingBlockExporter    bool   `mapstructure:"tracing_block_exporter"`
}

// ShutdownFunc shuts down a tracing env.
type ShutdownFunc func() error

// launcher launches a tracing env.
type launcher struct {
	// Exporter is the endpoint to which traces are exported.
	exporter *otlp.Exporter
	// Resource describes an application/service
	resource *resource.Resource
	// SpanProcessors process spans before getting exported. (pipeline pattern)
	spanprocs []trace.SpanProcessor
	// TracerProvider provides a tracer to intiate the starting of a span.
	provider *trace.TracerProvider
}
