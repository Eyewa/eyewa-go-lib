package exporter

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Mock mocks an otlptrace.Exporter
type Mock struct {
	mock.Mock
}

// NewMockExporter constructs a new mock exporter.
func NewMock() Exporter {
	return new(Mock)
}

func (mock *Mock) Start(ctx context.Context) error {
	args := mock.Called(ctx)
	return args.Error(0)
}

func (mock *Mock) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	args := mock.Called(ctx, spans)
	return args.Error(0)
}

func (mock *Mock) Shutdown(ctx context.Context) error {
	args := mock.Called(ctx)
	return args.Error(0)
}
