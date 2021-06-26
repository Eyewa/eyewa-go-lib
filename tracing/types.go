package tracing

import (
	"context"

	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

// Config represents the configuration used to configure a tracing environment.
type Config struct {
	ExporterEndpoint         string `mapstructure:"exporter_endpoint"`
	ExporterEndpointInsecure bool   `mapstructure:"exporter_endpoint_insecure"`
	ExporterBlocking         bool   `mapstructure:"exporter_blocking"`
	ServiceName              string `mapstructure:"service_name"`
	ServiceVersion           string `mapstructure:"service_version"`
	resourceAttributes       map[string]string
	Resource                 *resource.Resource
	logger                   *zap.Logger
}

// Exporter exports spans to a span destination.
type Exporter interface {
	// Start starts the exporter
	Start(ctx context.Context) error

	// ExportSpans exports spans to a span destination.
	ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error

	// Shutdown shuts down the connection to a span destination.
	Shutdown(ctx context.Context) error
}

// Option is a configuration option in a tracing environment.
type Option func(*Config)

// ShutdownFunc is the showdown function that
// shuts down a tracing environment.
type ShutdownFunc func() error

// Launcher is responsible for launching a
// tracing environment.
type Launcher struct {
	ctx      context.Context
	exporter Exporter
}
