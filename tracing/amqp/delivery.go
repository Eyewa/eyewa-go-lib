package amqp

// import (
// 	"context"
// 	"fmt"

// 	"github.com/streadway/amqp"
// 	"go.opentelemetry.io/otel/semconv"
// 	"go.opentelemetry.io/otel/trace"
// )

// var (
// 	deliverySpanName = fmt.Sprintf("%s.ConsumeDelivery", messagingSystem)
// )

// // StartDeliverySpan starts a span from an amqp.Delivery. It decorates
// // the span with amqp related attributes. If the amqp.Delivery
// // contains an existing trace context, then it will continue from there.
// func StartDeliverySpan(d amqp.Delivery, opts ...Option) (context.Context, EndSpanFunc, trace.Span) {
// 	cfg := newConfig(opts...)

// 	// set amqp message span attributes.
// 	spanOpts := []trace.SpanOption{
// 		trace.WithAttributes(
// 			semconv.MessagingSystemKey.String(messagingSystem),
// 			semconv.MessagingDestinationKindKeyQueue,
// 			semconv.MessagingOperationReceive,
// 			semconv.MessagingRabbitMQRoutingKeyKey.String(d.RoutingKey)),
// 		trace.WithSpanKind(trace.SpanKindConsumer),
// 	}

// 	// extract context from headers, if none, the context will use Background context.
// 	ctx := cfg.Propagators.Extract(context.Background(), HeaderCarrier(d.Headers))

// 	// start the span and and receive a new ctx containing the parent
// 	ctx, span := cfg.Tracer.Start(ctx, deliverySpanName, spanOpts...)

// 	// Inject current span context, so consumers can use it to propagate span.
// 	// cfg.Propagators.Inject(ctx, carrier)

// 	// return an end func to force the user of StartDeliverySpan to call end on the span.
// 	endFunc := span.End
// 	return ctx, endFunc, span
// }
