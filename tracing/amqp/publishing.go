package amqp

import (
	"context"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

type publishingSpan struct {
	publishing amqp.Publishing
	span       trace.Span
	cfg        config
}

func StartPublishingSpan(ctx context.Context, publishing amqp.Publishing, opts ...Option) (context.Context, func()) {
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
func (pubspan publishingSpan) start(ctx context.Context) (context.Context, trace.Span) {
	// If there's a span context in the message, use that as the parent context.
	carrier := NewPublishingCarrier(pubspan.publishing)
	ctx = pubspan.cfg.Propagators.Extract(ctx, carrier)

	// Create a span.
	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKey.String("rabbitmq"),
		semconv.MessagingDestinationKindKeyQueue,
	}

	// setup span options.
	opts := []trace.SpanOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindProducer),
	}

	// start the span and and receive a new ctx.
	ctx, span := pubspan.cfg.Tracer.Start(ctx, "rabbitmq.publish", opts...)

	// Inject new span context, so consumers can use it to propagate span.
	pubspan.cfg.Propagators.Inject(ctx, carrier)

	return ctx, span
}
