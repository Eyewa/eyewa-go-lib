package tracing

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv"
)

// constructs a new Resource with attributes.
func newResource() (*resource.Resource, error) {
	var attributes []attribute.KeyValue

	attributes = append(attributes,
		semconv.ServiceNameKey.String(config.ServiceName),
		semconv.HostNameKey.String(config.HostName),
	)

	// These detectors can't actually fail, ignoring the error.
	r, err := resource.New(
		context.Background(),
		resource.WithAttributes(attributes...),
	)

	return r, err
}
