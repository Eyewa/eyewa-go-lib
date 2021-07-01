package amqp

import (
	"context"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "github.com/eyewa/eyewa-go-lib/tracing/amqp"
const consumeSpanName = "rabbitmq.consume"
const messagingSystem = "rabbitmq"

type ConsumeFunc func(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)

// Wraps a channel consume function and starts a trace for every delivery.
func WrapConsume(consume ConsumeFunc) ConsumeFunc {
	return func(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
		// call the original consume function on the channel.
		deliveries, err := consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
		if err != nil {
			return nil, err
		}

		// start tracing and dispatching to the output channel.
		dispatcher := newDeliveryDispatcher(deliveries)
		go dispatcher.start()

		return dispatcher.out, nil
	}
}

type deliveryDispatcher struct {
	in  <-chan amqp.Delivery
	out chan amqp.Delivery
}

// newDeliveryDispatcher constructs a new delivery dispatcher.
func newDeliveryDispatcher(deliveries <-chan amqp.Delivery) *deliveryDispatcher {
	return &deliveryDispatcher{
		in:  deliveries,
		out: make(chan amqp.Delivery),
	}
}

// dispatch moves a delivery to the output channel.
func (dispatcher deliveryDispatcher) dispatch(delivery amqp.Delivery) {
	dispatcher.out <- delivery
}

// start starts tracing deliveries that come through the input channel
// and starts dispatching to the output.
func (dispatcher deliveryDispatcher) start() {
	for delivery := range dispatcher.in {
		endSpan := startTracing(delivery)
		dispatcher.dispatch(delivery)
		endSpan()
	}
	close(dispatcher.out)
}

// startTracing starts tracing a delivery
// and returns a function that ends the span.
func startTracing(delivery amqp.Delivery) func(...trace.SpanOption) {
	// Extract a span context from delivery.
	carrier := NewDeliveryHeaderCarrier(delivery)
	propagator := otel.GetTextMapPropagator()
	parentSpanContext := propagator.Extract(context.Background(), carrier)

	// Create a span.
	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKey.String(messagingSystem),
		semconv.MessagingDestinationKindKeyQueue,
		semconv.MessagingOperationReceive,
		semconv.MessagingMessageIDKey.String(delivery.MessageId),
		semconv.MessagingRabbitMQRoutingKeyKey.String(delivery.RoutingKey),
	}

	opts := []trace.SpanOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindConsumer),
	}

	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	newCtx, span := tracer.Start(parentSpanContext, consumeSpanName, opts...)

	// Inject current span context, so consumers can use it to propagate span.
	propagator.Inject(newCtx, carrier)

	return span.End
}
