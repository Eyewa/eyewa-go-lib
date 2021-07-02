package amqp

import (
	"context"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

func startProducerSpan(ctx context.Context, cfg config, publishing amqp.Publishing) trace.Span {
	// If there's a span context in the message, use that as the parent context.
	carrier := NewPublishingHeaderCarrier(&publishing)
	ctx = cfg.Propagators.Extract(ctx, carrier)

	// Create a span.
	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKey.String("rabbitmq"),
		semconv.MessagingDestinationKindKeyQueue,
	}
	opts := []trace.SpanOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindProducer),
	}
	ctx, span := cfg.Tracer.Start(ctx, "rabbitmq.publish", opts...)

	// Inject current span context, so consumers can use it to propagate span.
	cfg.Propagators.Inject(ctx, carrier)

	return span
}

func finishPublishingSpan(span trace.Span) {
	// span.SetAttributes(
	// 	semconv.MessagingMessageIDKey.String(strconv.FormatInt(offset, 10)),
	// 	kafkaPartitionKey.Int64(int64(partition)),
	// )
	// if err != nil {
	// 	span.SetStatus(codes.Error, err.Error())
	// }
	span.End()
}
