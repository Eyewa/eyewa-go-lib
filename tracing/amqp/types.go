package amqp

import (
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// deliverySpan is the span responsible for tracing an amqp.Delivery
// and attaching amqp related attributes.
type deliverySpan struct {
	cfg      config
	delivery amqp.Delivery
}

// deliverySpan is the span responsible for tracing an amqp.Publishing
// and attaching amqp related attributes.
type publishingSpan struct {
	publishing amqp.Publishing
	cfg        config
}

// config is used to configure starting of a trace
// on an amqp.Publishing and amqp.Delivery
type config struct {
	TracerProvider trace.TracerProvider
	Propagators    propagation.TextMapPropagator

	Tracer trace.Tracer
}

// Option is an amqp optional config
type Option func(*config)
