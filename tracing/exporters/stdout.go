package exporters

import (
	"context"

	liberrs "github.com/eyewa/eyewa-go-lib/errors"
	"github.com/eyewa/eyewa-go-lib/tracing"
	stdout "go.opentelemetry.io/otel/exporters/stdout"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func NewStdOut() (tracing.Exporter, error) {

	stdexp, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		return nil, liberrs.ErrExporterStartupFailure
	}
	exp := &stdOutExporter{exporter: stdexp}
	return exp, nil
}

// MethodNotImplemented
func (exp *stdOutExporter) Start(ctx context.Context) error {
	return nil
}

// MethodNotImplemented
func (exp *stdOutExporter) Shutdown(ctx context.Context) error {
	return nil
}

// MethodNotImplemented
func (exp *stdOutExporter) ExportSpans(ctx context.Context, spans []tracesdk.ReadOnlySpan) error {
	return nil
}
