package exporters

import (
	"context"

	"github.com/stretchr/testify/mock"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type MockOtelExporter struct {
	mock.Mock
}

func (mock *MockOtelExporter) Start(ctx context.Context) error {
	args := mock.Called(ctx)
	return args.Error(0)
}

func (mock *MockOtelExporter) Shutdown(ctx context.Context) error {
	args := mock.Called(ctx)
	return args.Error(0)
}

func (mock *MockOtelExporter) ExportSpans(ctx context.Context, spans []tracesdk.ReadOnlySpan) error {
	args := mock.Called(ctx, spans)
	return args.Error(0)
}
