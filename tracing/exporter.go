package tracing

import (
	"fmt"
	"time"

	"github.com/eyewa/eyewa-go-lib/log"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// constructs a new Open Telemetry Exporter.
func newOtelCollectorExporter() (*otlp.Exporter, error) {
	// configures exporter secure option
	var secureOpt otlpgrpc.Option
	if cfg.TracingSecureExporter {
		secureOpt = otlpgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOpt = otlpgrpc.WithInsecure()
	}
	log.Debug(fmt.Sprintf("Setting secure option for tracing exporter: %v", cfg.TracingSecureExporter))

	// configures exporter blocking option
	var blockingOpt otlpgrpc.Option
	if cfg.TracingBlockExporter {
		blockingOpt = otlpgrpc.WithDialOption(
			grpc.WithTimeout(2*time.Second),
			grpc.WithBlock(),
		)
	} else {
		blockingOpt = otlpgrpc.WithDialOption()
	}
	log.Debug(fmt.Sprintf("Setting blocking option for tracing exporter: %v", cfg.TracingBlockExporter))

	return otlp.NewUnstartedExporter(
		otlpgrpc.NewDriver(
			secureOpt,
			blockingOpt,
			otlpgrpc.WithEndpoint(cfg.TracingExporterEndpoint),
		),
	), nil
}
