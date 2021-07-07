package amqp

import (
	"context"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/oteltest"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

func TestStartPublishingSpan(t *testing.T) {
	propagators := propagation.TraceContext{}
	// var err error

	// setup a span recorder to store all spans
	sr := new(oteltest.SpanRecorder)
	provider := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))

	// setup an upstream span so that we have a context with a span and trace id to start with.
	parentCtx, _ := provider.Tracer(instrumentationName).Start(context.Background(), "test")

	// Create the publishing and start tracing to send spans to span recorder.
	msg := amqp.Publishing{Headers: amqp.Table{"test": "test"}}
	propagators.Inject(parentCtx, NewPublishingCarrier(&msg))
	ctx, span := StartPublishingSpan(parentCtx, &msg,
		WithTracerProvider(provider),
		WithPropagators(propagators),
	)
	span.End()

	spanList := sr.Completed()
	// Expected
	expectedList := []struct {
		attributeList []attribute.KeyValue
		parentSpanID  trace.SpanID
		kind          trace.SpanKind
	}{
		{
			attributeList: []attribute.KeyValue{
				semconv.MessagingSystemKey.String("rabbitmq"),
				semconv.MessagingDestinationKindKeyQueue,
			},
			parentSpanID: trace.SpanContextFromContext(ctx).SpanID(),
			kind:         trace.SpanKindProducer,
		},
	}

	for i, expected := range expectedList {
		span := spanList[i]

		// Check span
		assert.True(t, span.SpanContext().IsValid())
		assert.Equal(t, expected.kind, span.SpanKind())
		for _, k := range expected.attributeList {
			assert.Equal(t, k.Value, span.Attributes()[k.Key], k.Key)
		}

		// Check tracing propagation
		remoteSpanFromMessage := trace.SpanContextFromContext(
			propagators.Extract(context.Background(),
				NewPublishingCarrier(&msg)),
		)
		assert.True(t, remoteSpanFromMessage.IsValid())
	}

}
