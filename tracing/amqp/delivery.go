package amqp

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

var (
	instrumentationName = "github.com/eyewa/eyewa-go-lib/tracing/amqp"
	messagingSystem     = "rabbitmq"
	consumeSpanName     = fmt.Sprintf("%s.consume", messagingSystem)
	publishSpanName     = fmt.Sprintf("%s.publish", messagingSystem)
)

// deliverySpan is the span responsible for tracing an amqp.Delivery
// and attaching amqp related attributes.
type deliverySpan struct {
	cfg      config
	delivery amqp.Delivery
}

// StartDeliverySpan starts tracing a delivery and returns the new context and end span function.
func StartDeliverySpan(ctx context.Context, d amqp.Delivery, opts ...Option) (context.Context, func()) {
	cfg := newConfig(opts...)
	dspan := &deliverySpan{
		cfg:      cfg,
		delivery: d,
	}

	ctx, span := dspan.start(ctx)
	return ctx, func() {
		span.End()
	}
}

// start starts a span a delivery
// and returns a function that ends the span.
func (dspan deliverySpan) start(ctx context.Context) (context.Context, trace.Span) {
	// Extract a span context from delivery.
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

	tracer := dspan.cfg.TracerProvider.Tracer(instrumentationName)
	newCtx, span := tracer.Start(parentSpanContext, consumeSpanName, opts...)

	// Inject current span context, so consumers can use it to propagate span.
	dspan.cfg.Propagators.Inject(newCtx, carrier)

	return newCtx, span
}
