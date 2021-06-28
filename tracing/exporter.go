package tracing

import (
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// constructs a new Open Telemetry Exporter.
func newOtelCollectorExporter(endpoint string, secure, blocking bool) *otlp.Exporter {
	// secure connection by default
	secureOpt := otlpgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if !secure {
		secureOpt = otlpgrpc.WithInsecure()
	}
	// non blocking connection by default
	blockingOpt := otlpgrpc.WithDialOption()
	if blocking {
		blockingOpt = otlpgrpc.WithDialOption(grpc.WithBlock())
	}
	return otlp.NewUnstartedExporter(
		otlpgrpc.NewDriver(
			secureOpt,
			blockingOpt,
			otlpgrpc.WithEndpoint(endpoint),
		),
	)
}
