package rabbitmq

import (
	"context"

	"go.opentelemetry.io/otel/unit"

	"github.com/eyewa/eyewa-go-lib/errors"
	"github.com/eyewa/eyewa-go-lib/log"
	"github.com/eyewa/eyewa-go-lib/metrics"
	"go.opentelemetry.io/otel/metric"
)

type RabbitMQMetrics struct {
	PublishedEventCounter        *metrics.Counter
	FailedPublishedEventCounter  *metrics.Counter
	ConsumedEventCounter         *metrics.Counter
	FailedConsumedEventCounter   *metrics.Counter
	ActiveConsumingEventCounter  *metrics.UpDownCounter
	ConsumedEventLatencyRecorder *metrics.ValueRecorder
}

func NewRabbitMQMetrics() *RabbitMQMetrics {
	meter := metrics.NewMeter("rabbitmq.meter", context.Background())

	publishedEventCounter, err := meter.NewCounter("published.event.counter",
		metric.WithDescription("Counts published events"))
	if err != nil {
		log.Error(errors.ErrorFailedToCreateInstrument.Error())
	}

	failedPublishedEventCounter, err := meter.NewCounter("failed.published.event.counter",
		metric.WithDescription("Counts failed published events"))
	if err != nil {
		log.Error(errors.ErrorFailedToCreateInstrument.Error())
	}

	consumedEventCounter, err := meter.NewCounter("consumed.event.counter",
		metric.WithDescription("Counts consumed events"))
	if err != nil {
		log.Error(errors.ErrorFailedToCreateInstrument.Error())
	}

	failedConsumedEventCounter, err := meter.NewCounter("failed.consumed.event.counter",
		metric.WithDescription("Counts failed consumed events"))
	if err != nil {
		log.Error(errors.ErrorFailedToCreateInstrument.Error())
	}

	activeConsumingEventCounter, err := meter.NewUpDownCounter("active.consuming.event.counter",
		metric.WithDescription("Counts active consuming events"))
	if err != nil {
		log.Error(errors.ErrorFailedToCreateInstrument.Error())
	}

	consumedEventLatencyRecorder, err := meter.NewValueRecorder("consumed.event.latency.recorder",
		metric.WithUnit(unit.Milliseconds),
		metric.WithDescription("Records consumed event latency"))
	if err != nil {
		log.Error(errors.ErrorFailedToCreateInstrument.Error())
	}

	return &RabbitMQMetrics{
		PublishedEventCounter:        publishedEventCounter,
		FailedPublishedEventCounter:  failedPublishedEventCounter,
		ConsumedEventCounter:         consumedEventCounter,
		FailedConsumedEventCounter:   failedConsumedEventCounter,
		ActiveConsumingEventCounter:  activeConsumingEventCounter,
		ConsumedEventLatencyRecorder: consumedEventLatencyRecorder,
	}
}
