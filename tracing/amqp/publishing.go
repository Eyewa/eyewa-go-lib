package amqp

import (
	"context"
	"fmt"

	"github.com/eyewa/eyewa-go-lib/log"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

// StartDeliverySpan starts tracing a publishing and returns the new context and end span function.
func StartPublishingSpan(ctx context.Context, publishing *amqp.Publishing, opts ...Option) (context.Context, func()) {
	cfg := newConfig(opts...)
	pubspan := publishingSpan{
		publishing: publishing,
		cfg:        cfg,
	}

	ctx, span := pubspan.start(ctx)

	return ctx, func() {
		span.End()
	}
}

// start starts a span
// returns the new context and a function that ends the span.
func (pubspan publishingSpan) start(ctx context.Context) (context.Context, trace.Span) {
	// If there's a span context in the message, use that as the parent context.
	carrier := NewPublishingCarrier(pubspan.publishing)
	ctx = pubspan.cfg.Propagators.Extract(ctx, carrier)
	log.Info(fmt.Sprint(carrier.Keys()))

	// Create a span.
	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKey.String(messagingSystem),
		semconv.MessagingDestinationKindKeyQueue,
	}

	// setup span options.
	opts := []trace.SpanOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindProducer),
	}

	// start the span and and receive a new ctx.
	ctx, span := pubspan.cfg.Tracer.Start(ctx, publishSpanName, opts...)

	// Inject current span context, so publishers can use it to propagate span.
	pubspan.cfg.Propagators.Inject(ctx, carrier)

	return ctx, span
}
