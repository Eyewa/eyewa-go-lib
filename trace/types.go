package trace

// EnvConfig for tracing related environmental variables
type EnvConfig struct {
	// Exporter credentials
	CollectorEndpoint string `mapstructure:"trace_collector_endpoint"`
	// Name of the host
	HostName string `mapstructure:"host_name"`
	// Purely for identifying what service/service instance is being utilized by the resource
	ServiceName string `mapstructure:"service_name"`
}
