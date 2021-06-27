package exporter

import (
	"context"

	liberrs "github.com/eyewa/eyewa-go-lib/errors"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func NewStdOut(opts ...stdouttrace.Option) (Exporter, error) {
	opts = append(opts, stdouttrace.WithPrettyPrint())
	stdexp, err := stdouttrace.New(opts...)
	if err != nil {
		return nil, liberrs.ErrorExporterStartupFailure
	}
	exp := &stdOutExporter{exporter: stdexp}
	return exp, nil
}

func (exp *stdOutExporter) Start(ctx context.Context) error {
	// stdout has no start
	return nil
}

func (exp *stdOutExporter) Shutdown(ctx context.Context) error {
	err := exp.exporter.Shutdown(ctx)
	return liberrs.Wrap(err, liberrs.ErrorExporterShutdownFailure)
}

func (exp *stdOutExporter) ExportSpans(ctx context.Context, spans []tracesdk.ReadOnlySpan) error {
	err := exp.exporter.ExportSpans(ctx, spans)
	return liberrs.Wrap(err, liberrs.ErrorExporterShutdownFailure)
}
