package amqp

// StartDeliverySpan starts tracing a span and returns the end function.
// func StartPublishingSpan(ctx context.Context, publishing amqp.Publishing) (context.Context, func()) {
// 	ctx, endSpan := startTracing(ctx, publishing)
// 	return ctx, func() { endSpan() }
// }

// // startTracing starts tracing a publishing
// // and returns a function that ends the span.
// func startTracing(ctx context.Context, publishing amqp.Publishing) (context.Context, func(...trace.SpanOption)) {
// 	// Extract a span context from delivery.
// 	carrier := NewPublishingHeaderCarrier(&publishing)
// 	propagator := otel.GetTextMapPropagator()
// 	parentSpanContext := propagator.Extract(ctx, carrier)

// 	// Create a span.
// 	attrs := []attribute.KeyValue{
// 		semconv.MessagingSystemKey.String(messagingSystem),
// 		semconv.MessagingDestinationKindKeyQueue,
// 		semconv.MessagingOperationReceive,
// 		semconv.MessagingMessageIDKey.String(delivery.MessageId),
// 		semconv.MessagingRabbitMQRoutingKeyKey.String(delivery.RoutingKey),
// 	}

// 	opts := []trace.SpanOption{
// 		trace.WithAttributes(attrs...),
// 		trace.WithSpanKind(trace.SpanKindConsumer),
// 	}

// 	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
// 	newCtx, span := tracer.Start(parentSpanContext, consumeSpanName, opts...)

// 	// Inject current span context, so consumers can use it to propagate span.
// 	propagator.Inject(newCtx, carrier)

// 	return newCtx, span.End
// }
