package tracing

import (
	"time"

	"github.com/eyewa/eyewa-go-lib/errors"
	"github.com/eyewa/eyewa-go-lib/log"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// constructs a new Open Telemetry Exporter.
func newOtelCollectorExporter(endpoint string, secure, blocking bool) (*otlp.Exporter, error) {
	if endpoint == "" {
		return nil, errors.ErrorNoExporterEndpointSpecified
	}
	secureOpt := otlpgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if !secure {
		log.Debug("Setting insecure connection for tracing exporter.")
		secureOpt = otlpgrpc.WithInsecure()
	}

	blockingOpt := otlpgrpc.WithDialOption()
	if blocking {
		log.Debug("Setting blocking connection for tracing exporter.")
		blockingOpt = otlpgrpc.WithDialOption(
			grpc.WithTimeout(5*time.Second),
			grpc.WithBlock(),
		)
	}

	return otlp.NewUnstartedExporter(
		otlpgrpc.NewDriver(
			secureOpt,
			blockingOpt,
			otlpgrpc.WithEndpoint(endpoint),
		),
	), nil
}
