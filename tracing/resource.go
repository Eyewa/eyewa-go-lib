package tracing

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv"
)

// constructs a new Resource with attributes.
func newResource(ctx context.Context) (*resource.Resource, error) {
	var attributes []attribute.KeyValue

	attributes = append(attributes,
		semconv.ServiceNameKey.String(config.ServiceName),
		semconv.HostNameKey.String(config.HostName),
	)

	r, err := resource.New(ctx,
		resource.WithAttributes(attributes...),
	)

	return r, err
}
