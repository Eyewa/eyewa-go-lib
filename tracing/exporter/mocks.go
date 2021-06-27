package exporter

import (
	"context"

	"github.com/stretchr/testify/mock"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Mock mocks an otlptrace.Exporter
type StdOutMock struct {
	mock.Mock
}

// NewMockExporter constructs a new mock exporter.
func NewMock() Exporter {
	return new(StdOutMock)
}

func (mock *StdOutMock) Start(ctx context.Context) error {
	args := mock.Called(ctx)
	return args.Error(0)
}

func (mock *StdOutMock) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	args := mock.Called(ctx, spans)
	return args.Error(0)
}

func (mock *StdOutMock) Shutdown(ctx context.Context) error {
	args := mock.Called(ctx)
	return args.Error(0)
}

// Mock mocks an otlptrace.Exporter
type OtelColMock struct {
	mock.Mock
}

// NewOtelColMock constructs a new OtelColMock exporter.
func NewOtelColMock() Exporter {
	return new(StdOutMock)
}

func (mock *OtelColMock) Start(ctx context.Context) error {
	args := mock.Called(ctx)
	return args.Error(0)
}

func (mock *OtelColMock) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	args := mock.Called(ctx, spans)
	return args.Error(0)
}

func (mock *OtelColMock) Shutdown(ctx context.Context) error {
	args := mock.Called(ctx)
	return args.Error(0)
}

type MockWriter struct {
	mock.Mock
}

func (mock *MockWriter) Write(b []byte) (n int, err error) {
	args := mock.Called(b)
	return args.Int(0), args.Error(1)
}
