package tracing

import (
	"context"
	"os"

	"github.com/eyewa/eyewa-go-lib/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv"
)

// constructs a new Resource with attributes.
func newResource(svcName, svcVersion string) *resource.Resource {
	var attributes []attribute.KeyValue

	if len(svcName) > 0 {
		attributes = append(attributes, semconv.ServiceNameKey.String(svcName))
	}

	if len(svcVersion) > 0 {
		attributes = append(attributes, semconv.ServiceVersionKey.String(svcVersion))
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Debug("Failed to retrieve the hostname from the kernel.")
	} else {
		attributes = append(attributes, semconv.HostNameKey.String(hostname))
	}

	// These detectors can't actually fail, ignoring the error.
	r, _ := resource.New(
		context.Background(),
		resource.WithAttributes(attributes...),
	)

	return r
}