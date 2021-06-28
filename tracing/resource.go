package tracing

import (
	"context"
	"os"

	"github.com/eyewa/eyewa-go-lib/errors"

	"github.com/eyewa/eyewa-go-lib/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv"
)

// constructs a new Resource with attributes.
func newResource() (*resource.Resource, error) {
	var attributes []attribute.KeyValue

	if cfg.ServiceName == "" {
		return nil, errors.ErrorNoServiceNameSpecified
	}
	attributes = append(attributes, semconv.ServiceNameKey.String(cfg.ServiceName))

	// check if we can pickup the hostname from the os.
	hostname, err := os.Hostname()
	if err != nil {
		log.Debug("Failed to retrieve the hostname from the kernel.")
	} else {
		attributes = append(attributes, semconv.HostNameKey.String(hostname))
	}

	// These detectors can't actually fail, ignoring the error.
	r, err := resource.New(
		context.Background(),
		resource.WithAttributes(attributes...),
	)

	return r, err
}
