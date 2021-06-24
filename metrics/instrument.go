package metrics

import (
	"go.opentelemetry.io/otel/attribute"
)

// Add adds the value to the counter's sum. The labels should contain
// the keys and values to be associated with this value.
func (c Counter) Add(value float64, labels ...attribute.KeyValue) {
	c.Float64Counter.Add(c.ctx, value, labels...)
}

// Add adds the value to the counter's sum. The labels should contain
// the keys and values to be associated with this value.
func (u UpDownCounter) Add(value float64, labels ...attribute.KeyValue) {
	u.Float64UpDownCounter.Add(u.ctx, value, labels...)
}

// Record adds a new value to the list of ValueRecorder's records. The
// labels should contain the keys and values to be associated with
// this value.
func (v ValueRecorder) Record(value float64, labels ...attribute.KeyValue) {
	v.Float64ValueRecorder.Record(v.ctx, value, labels...)
}
