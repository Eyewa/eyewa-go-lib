package tracing

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/sdk/trace"
)

// MockExporter mocks the way an exporter behaves.
type MockExporter struct {
	mock.Mock
}

// NewMockExporter constructs a new mock exporter.
func NewMockExporter() Exporter {
	return new(MockExporter)
}

func (mock *MockExporter) Start(ctx context.Context) error {
	args := mock.Called(ctx)
	return args.Error(0)
}

func (mock *MockExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	args := mock.Called(ctx, spans)
	return args.Error(0)
}

func (mock *MockExporter) Shutdown(ctx context.Context) error {
	args := mock.Called(ctx)
	return args.Error(0)
}
