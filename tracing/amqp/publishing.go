package amqp

// import (
// 	"context"
// 	"fmt"

// 	"github.com/eyewa/eyewa-go-lib/log"
// 	"github.com/streadway/amqp"
// 	"go.opentelemetry.io/otel/attribute"
// 	"go.opentelemetry.io/otel/semconv"
// 	"go.opentelemetry.io/otel/trace"
// )

// var (
// 	publishingSpanName = fmt.Sprintf("%s.PublishPublishing", messagingSystem)
// )

// // StartDeliverySpan starts tracing a publishing and returns the new context and end span function.
// func StartPublishingSpan(publishing *amqp.Publishing, opts ...Option) (context.Context, trace.Span) {
// 	cfg := newConfig(opts...)
// 	pubspan := publishingSpan{
// 		publishing: publishing,
// 		cfg:        cfg,
// 	}

// 	ctx, span := pubspan.start()

// 	return ctx, span
// }

// // start starts a span
// // returns the new context and a function that ends the span.
// func (pubspan publishingSpan) start() (context.Context, trace.Span) {
// 	// If there's a span context in the message, use that as the parent context.
// 	carrier := HeaderCarrier(pubspan.publishing.Headers)
// 	ctx := pubspan.cfg.Propagators.Extract(context.Background(), carrier)
// 	log.Info(fmt.Sprintf("publishing carrier keys: %s", carrier.Keys()))

// 	// Create a span.
// 	attrs := []attribute.KeyValue{
// 		semconv.MessagingSystemKey.String(messagingSystem),
// 		semconv.MessagingDestinationKindKeyQueue,
// 	}

// 	// setup span options.
// 	opts := []trace.SpanOption{
// 		trace.WithAttributes(attrs...),
// 		trace.WithSpanKind(trace.SpanKindProducer),
// 	}

// 	// start the span and and receive a new ctx.
// 	ctx, span := pubspan.cfg.Tracer.Start(ctx, publishingSpanName, opts...)

// 	// Inject current span context, so publishers can use it to propagate span.
// 	pubspan.cfg.Propagators.Inject(ctx, carrier)

// 	return ctx, span
// }
