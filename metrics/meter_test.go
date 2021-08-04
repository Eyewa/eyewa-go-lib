package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMeter(t *testing.T) {
	meter := NewMeter("test.meter", nil)
	assert.NotNil(t, meter.Meter)
	assert.NotNil(t, meter.ctx)
}

func TestMeterNewCounter(t *testing.T) {
	meter := NewMeter("test.meter", nil)
	counter, err := meter.NewCounter("test.counter")

	assert.NotNil(t, counter.Float64Counter)
	assert.NotNil(t, counter.ctx)
	assert.Nil(t, err)
}

func TestNewUpDownMeterCounter(t *testing.T) {
	meter := NewMeter("test.meter", nil)
	upDownCounter, err := meter.NewUpDownCounter("test.upDownCounter")

	assert.NotNil(t, upDownCounter.Float64UpDownCounter)
	assert.NotNil(t, upDownCounter.ctx)
	assert.Nil(t, err)
}

func TestNewValueMeterRecorder(t *testing.T) {
	meter := NewMeter("test.meter", nil)
	valueRecorder, err := meter.NewValueRecorder("test.valueRecorder")

	assert.NotNil(t, valueRecorder.Float64ValueRecorder)
	assert.NotNil(t, valueRecorder.ctx)
	assert.Nil(t, err)
}

func TestNewAsyncMeterCounter(t *testing.T) {
	meter := NewMeter("test.meter", nil)
	sumObserver, err := meter.NewAsyncCounter("test.sumObserver", nil)

	assert.NotNil(t, sumObserver)
	assert.Nil(t, err)
}

func TestNewAsyncUpDownMeterCounter(t *testing.T) {
	meter := NewMeter("test.meter", nil)
	upDownSumObserver, err := meter.NewAsyncUpDownCounter("test.upDownSumObserver", nil)

	assert.NotNil(t, upDownSumObserver)
	assert.Nil(t, err)
}

func TestNewAsyncValueMeterRecorder(t *testing.T) {
	meter := NewMeter("test.meter", nil)
	valueObserver, err := meter.NewAsyncValueRecorder("test.valueObserver", nil)

	assert.NotNil(t, valueObserver)
	assert.Nil(t, err)
}
