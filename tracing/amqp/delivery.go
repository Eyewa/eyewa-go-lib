package amqp

import (
	"context"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

// StartDeliverySpan starts tracing a delivery and returns the new context and end span function.
func StartDeliverySpan(ctx context.Context, d *amqp.Delivery, opts ...Option) (context.Context, func()) {
	cfg := newConfig(opts...)
	dspan := deliverySpan{
		cfg:      cfg,
		delivery: d,
	}

	ctx, span := dspan.start(ctx)
	return ctx, func() {
		span.End()
	}
}

// start starts a span
// returns the new context and a function that ends the span.
func (dspan deliverySpan) start(ctx context.Context) (context.Context, trace.Span) {
	// If there's a span context in the message, use that as the parent context.
	carrier := NewDeliveryCarrier(dspan.delivery)
	parentSpanContext := dspan.cfg.Propagators.Extract(ctx, carrier)

	// Create a span.
	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKey.String(messagingSystem),
		semconv.MessagingDestinationKindKeyQueue,
		semconv.MessagingOperationReceive,
		semconv.MessagingMessageIDKey.String(dspan.delivery.MessageId),
		semconv.MessagingRabbitMQRoutingKeyKey.String(dspan.delivery.RoutingKey),
	}

	opts := []trace.SpanOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindConsumer),
	}

	// start the span and and receive a new ctx.
	newCtx, span := dspan.cfg.Tracer.Start(parentSpanContext, consumeSpanName, opts...)

	// Inject current span context, so consumers can use it to propagate span.
	dspan.cfg.Propagators.Inject(newCtx, carrier)

	return newCtx, span
}
