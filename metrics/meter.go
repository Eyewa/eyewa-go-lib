package metrics

import (
	"context"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
)

// NewMeter creates a new Meter
func NewMeter(name string, ctx context.Context, opts ...metric.MeterOption) *Meter {
	if ctx == nil {
		ctx = context.Background()
	}

	return &Meter{
		Meter: global.Meter(name, opts...),
		ctx:   ctx,
	}
}

// NewCounter creates a new Counter instrumentation for Meter
func (m *Meter) NewCounter(name string, iop ...metric.InstrumentOption) Counter {
	counter := metric.Must(m.Meter).NewFloat64Counter(name, iop...)

	return Counter{
		counter,
		m.ctx,
	}
}

// NewUpDownCounter creates a new UpDownCounter instrumentation for Meter
func (m *Meter) NewUpDownCounter(name string, iop ...metric.InstrumentOption) UpDownCounter {
	upDownCounter := metric.Must(m.Meter).NewFloat64UpDownCounter(name, iop...)

	return UpDownCounter{
		upDownCounter,
		m.ctx,
	}
}

// NewValueRecorder creates a new ValueRecorder instrumentation for Meter
func (m *Meter) NewValueRecorder(name string, iop ...metric.InstrumentOption) ValueRecorder {
	valueRecorder := metric.Must(m.Meter).NewFloat64ValueRecorder(name, iop...)

	return ValueRecorder{
		valueRecorder,
		m.ctx,
	}
}

// NewAsyncCounter creates a new AsyncCounter instrumentation for Meter
func (m *Meter) NewAsyncCounter(name string, cb Float64ObserverCallback, iop ...metric.InstrumentOption) AsyncCounter {
	sumObserver := metric.Must(m.Meter).NewFloat64SumObserver(name, metric.Float64ObserverFunc(cb), iop...)

	return AsyncCounter(sumObserver)
}

// NewAsyncUpDownCounter creates a new AsyncUpDownCounter for Meter
func (m *Meter) NewAsyncUpDownCounter(name string, cb Float64ObserverCallback, iop ...metric.InstrumentOption) AsyncUpDownCounter {
	upDownSumObserver := metric.Must(m.Meter).NewFloat64UpDownSumObserver(name, metric.Float64ObserverFunc(cb), iop...)

	return AsyncUpDownCounter(upDownSumObserver)
}

// NewAsyncValueRecorder creates a new AsyncValueRecorder for Meter
func (m *Meter) NewAsyncValueRecorder(name string, cb Float64ObserverCallback, iop ...metric.InstrumentOption) AsyncValueRecorder {
	valueObserver := metric.Must(m.Meter).NewFloat64ValueObserver(name, metric.Float64ObserverFunc(cb), iop...)

	return AsyncValueRecorder(valueObserver)
}
