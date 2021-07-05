package metrics

import (
	"context"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"time"
)

// Launcher is used for serving metrics.
type Launcher struct {
	exporter                *prometheus.Exporter
	enableHostInstrument    bool
	enableRuntimeInstrument bool
}

// Meter reports given instruments specified by OpenTelemetry
type Meter struct {
	metric.Meter

	ctx context.Context
}

// ExportOption is configuration for MetricLauncher.
type ExportOption struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`

	// CollectPeriod sets period interval of exporting process.
	CollectPeriod time.Duration `mapstructure:"METRIC_COLLECT_PERIOD"`
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

// AsyncCounter is a metrics that captures a precomputed sum of
// float64 values at a point in time.
type AsyncCounter metric.Float64SumObserver

// AsyncUpDownCounter is a metrics that captures a precomputed sum of
// float64 values at a point in time.
type AsyncUpDownCounter metric.Float64UpDownSumObserver

// AsyncValueRecorder is a metrics that captures a set of float64 values
// at a point in time.
type AsyncValueRecorder metric.Float64ValueObserver

// Float64ObserverCallback is for asynchronous metrics.
// The Callback function reports the absolute value of the counter
// User code is recommended not to provide more than one Measurement with the same attributes in a single callback.
// If it happens, the SDK can decide how to handle it. For example, during the callback invocation
// if two measurements value=1, attributes={pid:4, bitness:64} and value=2, attributes={pid:4, bitness:64} are reported,
// the SDK can decide to simply let them pass through (so the downstream consumer can handle duplication),
// drop the entire data, pick the last one, or something else.
// The API must treat observations from a single callback as logically taking place at a single instant,
// such that when recorded, observations from a single callback MUST be reported with identical timestamps.
type Float64ObserverCallback metric.Float64ObserverFunc
