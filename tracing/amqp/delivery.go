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

// StartDeliverySpan starts tracing a span and returns the end function.
func StartDeliverySpan(ctx context.Context, delivery amqp.Delivery) (context.Context, func()) {
	ctx, endSpan := startTracing(ctx, delivery)
	return ctx, func() { endSpan() }
}

// startTracing starts tracing a delivery
// and returns a function that ends the span.
func startTracing(ctx context.Context, delivery amqp.Delivery) (context.Context, func(...trace.SpanOption)) {
	// Extract a span context from delivery.
	carrier := NewDeliveryHeaderCarrier(delivery)
	propagator := otel.GetTextMapPropagator()
	parentSpanContext := propagator.Extract(ctx, carrier)

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

	return newCtx, span.End
}
