package tracing

import (
	"fmt"

	"github.com/eyewa/eyewa-go-lib/log"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// constructs a new Open Telemetry Exporter.
func newOtelExporter() (*otlp.Exporter, error) {
	// configures exporter secure option
	var secureOpt otlpgrpc.Option
	if config.TracingSecureExporter {
		secureOpt = otlpgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOpt = otlpgrpc.WithInsecure()
	}
	log.Debug(fmt.Sprintf("Setting secure option for tracing exporter: %v", config.TracingSecureExporter))

	// configures exporter blocking option
	var blockingOpt otlpgrpc.Option
	if config.TracingBlockExporter {
		blockingOpt = otlpgrpc.WithDialOption(
			grpc.WithTimeout(exporterTimeout),
			grpc.WithBlock(),
		)
	} else {
		blockingOpt = otlpgrpc.WithDialOption()
	}
	log.Debug(fmt.Sprintf("Setting blocking option for tracing exporter: %v", config.TracingBlockExporter))

	return otlp.NewUnstartedExporter(
		otlpgrpc.NewDriver(
			secureOpt,
			blockingOpt,
			otlpgrpc.WithEndpoint(config.TracingExporterEndpoint),
		),
	), nil
}
