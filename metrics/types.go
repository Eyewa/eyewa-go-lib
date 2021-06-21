package metrics

import (
	"context"
	"go.opentelemetry.io/otel/metric"
)

// Meter reports given instruments specified by OpenTelemetry
type Meter struct {
	metric.Meter

	ctx context.Context
}

// Counter is counter instrument accumulates float64 values
type Counter struct {
	metric.Float64Counter

	ctx context.Context
}

// UpDownCounter is a metrics instrument that sums floating
// point values.
type UpDownCounter struct {
	metric.Float64UpDownCounter

	ctx context.Context
}

// ValueRecorder is a metrics that records float64 values.
type ValueRecorder struct {
	metric.Float64ValueRecorder

	ctx context.Context
}

// SumObserver is a metrics that captures a precomputed sum of
// float64 values at a point in time.
type SumObserver metric.Float64SumObserver

// UpDownSumObserver is a metrics that captures a precomputed sum of
// float64 values at a point in time.
type UpDownSumObserver metric.Float64UpDownSumObserver

// ValueObserver is a metrics that captures a set of float64 values
// at a point in time.
type ValueObserver metric.Float64ValueObserver

// Callback is for asynchronous metrics.
// The Callback function reports the absolute value of the counter
// User code is recommended not to provide more than one Measurement with the same attributes in a single callback.
// If it happens, the SDK can decide how to handle it. For example, during the callback invocation
// if two measurements value=1, attributes={pid:4, bitness:64} and value=2, attributes={pid:4, bitness:64} are reported,
// the SDK can decide to simply let them pass through (so the downstream consumer can handle duplication),
// drop the entire data, pick the last one, or something else.
// The API must treat observations from a single callback as logically taking place at a single instant,
// such that when recorded, observations from a single callback MUST be reported with identical timestamps.
type MetricsCallback metric.Float64ObserverFunc