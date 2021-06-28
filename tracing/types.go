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
	config    config
	exporter  *otlp.Exporter
	resource  *resource.Resource
	spanprocs []trace.SpanProcessor
	provider  *trace.TracerProvider
}
