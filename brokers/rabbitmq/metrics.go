package rabbitmq

import (
	"context"

	"go.opentelemetry.io/otel/unit"

	"github.com/eyewa/eyewa-go-lib/errors"
	"github.com/eyewa/eyewa-go-lib/log"
	"github.com/eyewa/eyewa-go-lib/metrics"
	"go.opentelemetry.io/otel/metric"
)

// RabbitMQMetrics is a collection of standart metrics
type RabbitMQMetrics struct {
	PublishedEventCounter           *metrics.Counter
	PublishEventFailureCounter      *metrics.Counter
	ConsumedEventCounter            *metrics.Counter
	UnmarshalEventFailureCounter    *metrics.Counter
	MarshalEventFailureCounter      *metrics.Counter
	NackFailureCounter              *metrics.Counter
	DeadletterPublishFailureCounter *metrics.Counter
	ActiveConsumingEventCounter     *metrics.UpDownCounter
	ConsumedEventLatencyRecorder    *metrics.ValueRecorder
}

// NewRabbitMQMetrics creates a instance of RabbitMQMetrics
func NewRabbitMQMetrics() *RabbitMQMetrics {
	meter := metrics.NewMeter("rabbitmq.meter", context.Background())

	publishedEventCounter, err := meter.NewCounter("published.event.counter",
		metric.WithDescription("Counts published events"))
	if err != nil {
		log.Error(errors.ErrorFailedToCreateInstrument.Error())
	}

	publishEventFailureCounter, err := meter.NewCounter("publish.event.failure.counter",
		metric.WithDescription("Counts failed published events"))
	if err != nil {
		log.Error(errors.ErrorFailedToCreateInstrument.Error())
	}

	consumedEventCounter, err := meter.NewCounter("consumed.event.counter",
		metric.WithDescription("Counts consumed events"))
	if err != nil {
		log.Error(errors.ErrorFailedToCreateInstrument.Error())
	}

	marshalEventFailureCounter, err := meter.NewCounter("marshal.event.failure.counter",
		metric.WithDescription("Counts marshal event failures"))
	if err != nil {
		log.Error(errors.ErrorFailedToCreateInstrument.Error())
	}

	unmarshalEventFailureCounter, err := meter.NewCounter("unmarshal.event.failure.counter",
		metric.WithDescription("Counts unmarshal event failures"))
	if err != nil {
		log.Error(errors.ErrorFailedToCreateInstrument.Error())
	}

	nackFailureCounter, err := meter.NewCounter("nack.failure.counter",
		metric.WithDescription("Counts nack failures"))
	if err != nil {
		log.Error(errors.ErrorFailedToCreateInstrument.Error())
	}

	deadletterPublishFailureCounter, err := meter.NewCounter("deadletter.publish.failure.counter",
		metric.WithDescription("Counts deadletter publishing failures"))
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
		PublishedEventCounter:           publishedEventCounter,
		PublishEventFailureCounter:      publishEventFailureCounter,
		ConsumedEventCounter:            consumedEventCounter,
		MarshalEventFailureCounter:      marshalEventFailureCounter,
		UnmarshalEventFailureCounter:    unmarshalEventFailureCounter,
		NackFailureCounter:              nackFailureCounter,
		DeadletterPublishFailureCounter: deadletterPublishFailureCounter,
		ActiveConsumingEventCounter:     activeConsumingEventCounter,
		ConsumedEventLatencyRecorder:    consumedEventLatencyRecorder,
	}
}
