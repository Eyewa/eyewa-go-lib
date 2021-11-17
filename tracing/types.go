package tracing

const ProductETLInstrumentationName string = "github.com/eyewa/product-etl-service"

// Config is the tracing environment configuration.
type Config struct {
	// Required
	ServiceName             string `mapstructure:"service_name"`
	TracingExporterEndpoint string `mapstructure:"tracing_exporter_endpoint"`

	// Optional
	TracingSecureExporter bool   `mapstructure:"tracing_secure_exporter"`
	TracingBlockExporter  bool   `mapstructure:"tracing_block_exporter"`
	HostName              string `mapstructure:"hostname"`
}

// ShutdownFunc shuts down a tracing env.
type ShutdownFunc func() error
