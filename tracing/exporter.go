package tracing

import (
	"context"
	"fmt"

	"github.com/eyewa/eyewa-go-lib/log"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// constructs a new Open Telemetry Exporter.
func newOtelExporter(ctx context.Context) (*otlp.Exporter, error) {
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

	// start the exporter
	exporter, err := otlp.NewExporter(ctx,
		otlpgrpc.NewDriver(
			secureOpt,
			blockingOpt,
			otlpgrpc.WithEndpoint(config.TracingExporterEndpoint),
		),
	)

	if err != nil {
		return nil, err
	}

	return exporter, nil
}
