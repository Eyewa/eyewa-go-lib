package metrics

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv"
)

// newResource constructs a new Resource with attributes.
func newResource(option exportOption) (*resource.Resource, error) {
	var attributes []attribute.KeyValue

	attributes = append(attributes,
		semconv.ServiceNameKey.String(option.serviceName),
	)

	// These detectors can't actually fail, ignoring the error.
	r, err := resource.New(
		context.Background(),
		resource.WithAttributes(attributes...),
	)

	return r, err
}
