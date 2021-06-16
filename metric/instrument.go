package metric

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Counter is counter instrument accumulates float64 values
type Counter struct {
	metric.Float64Counter

	ctx context.Context
}

// Add adds the value to the counter's sum. The labels should contain
// the keys and values to be associated with this value.
func (b Counter) Add(value float64, labels ...attribute.KeyValue) {
	b.Float64Counter.Add(b.ctx, value, labels...)
}

// UpDownCounter is a metric instrument that sums floating
// point values.
type UpDownCounter struct {
	metric.Float64UpDownCounter

	ctx context.Context
}

// Add adds the value to the counter's sum. The labels should contain
// the keys and values to be associated with this value.
func (u UpDownCounter) Add(value float64, labels ...attribute.KeyValue) {
	u.Float64UpDownCounter.Add(u.ctx, value, labels...)
}

// ValueRecorder is a metric that records float64 values.
type ValueRecorder struct {
	metric.Float64ValueRecorder

	ctx context.Context
}

// Record adds a new value to the list of ValueRecorder's records. The
// labels should contain the keys and values to be associated with
// this value.
func (v ValueRecorder) Record(value float64, labels ...attribute.KeyValue) {
	v.Float64ValueRecorder.Record(v.ctx, value, labels...)
}

// SumObserver is a metric that captures a precomputed sum of
// float64 values at a point in time.
type SumObserver metric.Float64SumObserver

// UpDownSumObserver is a metric that captures a precomputed sum of
// float64 values at a point in time.
type UpDownSumObserver metric.Float64UpDownSumObserver

// ValueObserver is a metric that captures a set of float64 values
// at a point in time.
type ValueObserver metric.Float64ValueObserver
