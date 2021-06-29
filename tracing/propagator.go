package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// registersPropagators registers the trace propagation format.
// https://opentelemetry.lightstep.com/core-concepts/context-propagation/
func registerPropagators() {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))
}
