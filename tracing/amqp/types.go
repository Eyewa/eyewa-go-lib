package amqp

import (
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// deliverySpan is the span responsible for tracing an amqp.Delivery
// and attaching amqp related attributes.
type deliverySpan struct {
	delivery *amqp.Delivery
	cfg      config
}

// deliverySpan is the span responsible for tracing an amqp.Publishing
// and attaching amqp related attributes.
type publishingSpan struct {
	publishing *amqp.Publishing
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

// DeliveryCarrier injects and extracts
// traces from the headers of a amqp.Delivery.
type DeliveryCarrier struct {
	delivery *amqp.Delivery
}

// PublishingCarrier injects and extracts
// traces from the headers of a amqp.Publishing.
type PublishingCarrier struct {
	publishing *amqp.Publishing
}
