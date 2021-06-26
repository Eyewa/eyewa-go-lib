package exporters

import (
	"context"

	liberrs "github.com/eyewa/eyewa-go-lib/errors"
	"github.com/eyewa/eyewa-go-lib/tracing"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//
func NewOpenTelemetryCollectorExporter(endpoint string, blocking, insecure bool) (tracing.Exporter, error) {
	var opts []otlpgrpc.Option
	secureOption := otlpgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if insecure {
		secureOption = otlpgrpc.WithInsecure()
	}

	blockingOption := otlpgrpc.WithDialOption()
	if blocking {
		blockingOption = otlpgrpc.WithDialOption(grpc.WithBlock())
	}

	opts = append(opts,
		secureOption,
		blockingOption,
		otlpgrpc.WithEndpoint(endpoint),
	)

	exp, err := otlp.NewExporter(context.Background(), otlpgrpc.NewDriver(opts...))
	if err != nil {
		return nil, liberrs.ErrExporterStartupFailure
	}

	exporter := &otelCollectorExporter{exporter: exp}
	return exporter, nil
}

func (exp *otelCollectorExporter) Start(ctx context.Context) error {
	return exp.exporter.Start(ctx)
}

func (exp *otelCollectorExporter) Shutdown(ctx context.Context) error {
	return nil
}

func (exp *otelCollectorExporter) ExportSpans(ctx context.Context, spans []tracesdk.ReadOnlySpan) error {
	return nil
}
