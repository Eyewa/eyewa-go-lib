package tracing

import (
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// config is the tracing environment configuration.
type config struct {
	ServiceName      string `mapstructure:"service_name"`
	ServiceVersion   string `mapstructure:"service_version"`
	ExporterEndpoint string `mapstructure:"exporter_endpoint"`
	ExporterSecure   bool   `mapstructure:"exporter_secure"`
	ExporterBlocking bool   `mapstructure:"exporter_blocking"`
}

// ShutdownFunc shuts down a tracing env.
type ShutdownFunc func() error

// launcher launches a tracing env.
type launcher struct {
	started bool
	config  config
	// Exporter is the endpoint to which traces are exported.
	exporter *otlp.Exporter
	// Resource describes an application/service
	resource *resource.Resource
	// SpanProcessors process spans before getting exported. (pipeline pattern)
	spanprocs []trace.SpanProcessor
	// TracerProvider provides a tracer to intiate the starting of a span.
	provider *trace.TracerProvider
}
