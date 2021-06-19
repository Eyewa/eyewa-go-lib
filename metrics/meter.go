package metrics

import (
	"context"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
)

// Meter reports given instruments specified by OpenTelemetry
type Meter struct {
	metric.Meter

	ctx context.Context
}

// NewMeter creates a new Meter
func NewMeter(name string, ctx context.Context) *Meter {
	if ctx == nil {
		ctx = context.Background()
	}

	return &Meter{
		Meter: global.Meter(name,metric.WithInstrumentationVersion("1.0.0")),
		ctx:   ctx,
	}
}

// NewCounter creates new Counter instrumentation for Meter
func (m *Meter) NewCounter(name string, iop ...metric.InstrumentOption) Counter {
	counter := metric.Must(m.Meter).NewFloat64Counter(name, iop...)

	return Counter{
		counter,
		m.ctx,
	}
}

// NewUpDownCounter creates new UpDownCounter instrumentation for Meter
func (m *Meter) NewUpDownCounter(name string, iop ...metric.InstrumentOption) UpDownCounter {
	upDownCounter := metric.Must(m.Meter).NewFloat64UpDownCounter(name, iop...)

	return UpDownCounter{
		upDownCounter,
		m.ctx,
	}
}

// NewValueRecorder creates new ValueRecorder instrumentation for Meter
func (m *Meter) NewValueRecorder(name string, iop ...metric.InstrumentOption) ValueRecorder {
	valueRecorder := metric.Must(m.Meter).NewFloat64ValueRecorder(name, iop...)

	return ValueRecorder{
		valueRecorder,
		m.ctx,
	}
}

// Callback is for asynchronous metrics.
// The Callback function reports the absolute value of the counter
// User code is recommended not to provide more than one Measurement with the same attributes in a single callback.
// If it happens, the SDK can decide how to handle it. For example, during the callback invocation
// if two measurements value=1, attributes={pid:4, bitness:64} and value=2, attributes={pid:4, bitness:64} are reported,
// the SDK can decide to simply let them pass through (so the downstream consumer can handle duplication),
// drop the entire data, pick the last one, or something else.
// The API must treat observations from a single callback as logically taking place at a single instant,
// such that when recorded, observations from a single callback MUST be reported with identical timestamps.
type Callback metric.Float64ObserverFunc

// NewSumObserver creates new SumObserver instrumentation for Meter
func (m *Meter) NewSumObserver(name string, cb Callback, iop ...metric.InstrumentOption) SumObserver {
	sumObserver := metric.Must(m.Meter).NewFloat64SumObserver(name, metric.Float64ObserverFunc(cb), iop...)

	return SumObserver(sumObserver)
}

// NewUpDownSumObserver creates new UpDownSumObserver for Meter
func (m *Meter) NewUpDownSumObserver(name string, cb Callback, iop ...metric.InstrumentOption) UpDownSumObserver {
	upDownSumObserver := metric.Must(m.Meter).NewFloat64UpDownSumObserver(name, metric.Float64ObserverFunc(cb), iop...)

	return UpDownSumObserver(upDownSumObserver)
}

// NewValueObserver creates new ValueObserver for Meter
func (m *Meter) NewValueObserver(name string, cb Callback, iop ...metric.InstrumentOption) ValueObserver {
	valueObserver := metric.Must(m.Meter).NewFloat64ValueObserver(name, metric.Float64ObserverFunc(cb), iop...)

	return ValueObserver(valueObserver)
}
