package exporter

// import (
// 	"context"

// 	liberrs "github.com/eyewa/eyewa-go-lib/errors"
// 	"go.opentelemetry.io/otel/exporters/otlp"
// 	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
// 	tracesdk "go.opentelemetry.io/otel/sdk/trace"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials"
// )

// // NewOpenTelemetryCollectorExporter constructs an exporter that exporters to an open telemetry collector.
// func NewOpenTelemetryCollectorExporter(endpoint string, blocking, insecure bool) (Exporter, error) {
// 	var opts []otlpgrpc.Option
// 	secureOption := otlpgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
// 	if insecure {
// 		secureOption = otlpgrpc.WithInsecure()
// 	}

// 	blockingOption := otlpgrpc.WithDialOption()
// 	if blocking {
// 		blockingOption = otlpgrpc.WithDialOption(grpc.WithBlock())
// 	}

// 	opts = append(opts,
// 		secureOption,
// 		blockingOption,
// 		otlpgrpc.WithEndpoint(endpoint),
// 	)

// 	exp, err := otlp.NewExporter(context.Background(), otlpgrpc.NewDriver(opts...))
// 	if err != nil {
// 		return nil, liberrs.ErrExporterStartupFailure
// 	}

// 	exporter := &otelCollectorExporter{exporter: exp}
// 	return exporter, nil
// }

// func (exp *otelCollectorExporter) Start(ctx context.Context) error {
// 	err := exp.exporter.Start(ctx)
// 	return liberrs.Wrap(err, liberrs.ErrExporterShutdownFailure)
// }

// func (exp *otelCollectorExporter) Shutdown(ctx context.Context) error {
// 	err := exp.exporter.Shutdown(ctx)
// 	return liberrs.Wrap(err, liberrs.ErrExporterShutdownFailure)
// }

// func (exp *otelCollectorExporter) ExportSpans(ctx context.Context, spans []tracesdk.ReadOnlySpan) error {
// 	err := exp.exporter.ExportSpans(ctx)
// 	return liberrs.Wrap(err, liberrs.ErrExporterShutdownFailure)
// }
