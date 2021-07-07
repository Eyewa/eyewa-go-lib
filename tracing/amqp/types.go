package amqp

import (
	"github.com/streadway/amqp"
)

// deliverySpan wraps an amqp.Delivery for tracing.
// type delivery struct {
// 	*amqp.Delivery
// 	cfg config
// }

// // deliverySpan wraps an amqp.Publishing for tracing.
// type publishingSpan struct {
// 	publishing *amqp.Publishing
// 	cfg        config
// }

// config is used to configure starting of a trace
// on an amqp.Publishing and amqp.Delivery
// type config struct {
// 	TracerProvider trace.TracerProvider
// 	Propagators    propagation.TextMapPropagator

// 	Tracer trace.Tracer
// }

// // Option is an amqp optional config
// type Option func(*config)

// HeaderCarrier adapts amqp.Table to satisfy the TextMapCarrier interface.
type HeaderCarrier amqp.Table
