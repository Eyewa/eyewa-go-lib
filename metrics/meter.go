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
func (m *Meter) NewCounter(name string, iop ...metric.InstrumentOption) (*Counter, error) {
	counter, err := m.Meter.NewFloat64Counter(name, iop...)
	if err != nil {
		return nil, err
	}

	return &Counter{
		counter,
		m.ctx,
	}, nil
}

// NewUpDownCounter creates a new UpDownCounter instrumentation for Meter
func (m *Meter) NewUpDownCounter(name string, iop ...metric.InstrumentOption) (*UpDownCounter, error) {
	upDownCounter, err := m.Meter.NewFloat64UpDownCounter(name, iop...)
	if err != nil {
		return nil, err
	}

	return &UpDownCounter{
		upDownCounter,
		m.ctx,
	}, nil
}

// NewValueRecorder creates a new ValueRecorder instrumentation for Meter
func (m *Meter) NewValueRecorder(name string, iop ...metric.InstrumentOption) (*ValueRecorder, error) {
	valueRecorder, err := m.Meter.NewFloat64ValueRecorder(name, iop...)
	if err != nil {
		return nil, err
	}

	return &ValueRecorder{
		valueRecorder,
		m.ctx,
	}, nil
}

// NewAsyncCounter creates a new AsyncCounter instrumentation for Meter
func (m *Meter) NewAsyncCounter(name string, cb Float64ObserverCallback, iop ...metric.InstrumentOption) (*AsyncCounter, error) {
	sumObserver, err := m.Meter.NewFloat64SumObserver(name, metric.Float64ObserverFunc(cb), iop...)
	if err != nil {
		return nil, err
	}

	asyncCounter := AsyncCounter(sumObserver)
	return &asyncCounter, nil
}

// NewAsyncUpDownCounter creates a new AsyncUpDownCounter for Meter
func (m *Meter) NewAsyncUpDownCounter(name string, cb Float64ObserverCallback, iop ...metric.InstrumentOption) (*AsyncUpDownCounter, error) {
	upDownSumObserver, err := m.Meter.NewFloat64UpDownSumObserver(name, metric.Float64ObserverFunc(cb), iop...)
	if err != nil {
		return nil, err
	}

	asyncUpDownCounter := AsyncUpDownCounter(upDownSumObserver)
	return &asyncUpDownCounter, nil
}

// NewAsyncValueRecorder creates a new AsyncValueRecorder for Meter
func (m *Meter) NewAsyncValueRecorder(name string, cb Float64ObserverCallback, iop ...metric.InstrumentOption) (*AsyncValueRecorder, error) {
	valueObserver, err := m.Meter.NewFloat64ValueObserver(name, metric.Float64ObserverFunc(cb), iop...)
	if err != nil {
		return nil, err
	}

	asyncValueRecorder := AsyncValueRecorder(valueObserver)
	return &asyncValueRecorder, nil
}
