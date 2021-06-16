package metric

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMeter(t *testing.T) {
	var meter *Meter

	meter = NewMeter("test.meter", nil)
	assert.NotNil(t, meter.Meter)
	assert.NotNil(t, meter.ctx)
}

func TestMeter_NewCounter(t *testing.T) {
	meter := NewMeter("test.meter", nil)
	counter := meter.NewCounter("test.counter")

	assert.NotNil(t, counter.Float64Counter)
	assert.NotNil(t, counter.ctx)
}

func TestMeter_NewUpDownCounter(t *testing.T) {
	meter := NewMeter("test.meter", nil)
	upDownCounter := meter.NewUpDownCounter("test.upDownCounter")

	assert.NotNil(t, upDownCounter.Float64UpDownCounter)
	assert.NotNil(t, upDownCounter.ctx)
}

func TestMeter_NewValueRecorder(t *testing.T) {
	meter := NewMeter("test.meter", nil)
	valueRecorder := meter.NewValueRecorder("test.valueRecorder")

	assert.NotNil(t, valueRecorder.Float64ValueRecorder)
	assert.NotNil(t, valueRecorder.ctx)
}

func TestMeter_NewSumObserver(t *testing.T) {
	meter := NewMeter("test.meter", nil)
	sumObserver := meter.NewSumObserver("test.sumObserver", nil)

	assert.NotNil(t, sumObserver)
}

func TestMeter_NewUpDownSumObserver(t *testing.T) {
	meter := NewMeter("test.meter", nil)
	upDownSumObserver := meter.NewUpDownSumObserver("test.upDownSumObserver", nil)

	assert.NotNil(t, upDownSumObserver)
}

func TestMeter_NewValueObserver(t *testing.T) {
	meter := NewMeter("test.meter", nil)
	valueObserver := meter.NewValueObserver("test.valueObserver", nil)

	assert.NotNil(t, valueObserver)
}
