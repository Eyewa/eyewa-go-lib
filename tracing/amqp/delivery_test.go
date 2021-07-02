package amqp

import (
	"context"
	"fmt"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/oteltest"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

type expectedList struct {
	attributeList []attribute.KeyValue
	parentSpanID  trace.SpanID
	kind          trace.SpanKind
	msgKey        []byte
}

func injectTraceInfo(parentCtx context.Context, d amqp.Delivery) {
	propagator := propagation.TraceContext{}
	propagator.Inject(parentCtx, NewDeliveryCarrier(d))
}

func TestStartDeliverySpan(t *testing.T) {
	// setup a span recorder to store all spans
	sr := new(oteltest.SpanRecorder)
	provider := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))

	// setup an upstream span so that we have a context with a span and trace id.
	parentCtx, _ := provider.Tracer(instrumentationName).Start(context.Background(), "")

	tests := []struct {
		delivery amqp.Delivery
		expected expectedList
	}{
		{
			delivery: amqp.Delivery{
				Body:       []byte("this is a test"),
				Headers:    amqp.Table{"testkey": "testvalue"},
				RoutingKey: "routing.key",
			},
			expected: expectedList{
				attributeList: []attribute.KeyValue{
					semconv.MessagingSystemKey.String("rabbitmq"),
					semconv.MessagingDestinationKindKeyQueue,
					semconv.MessagingRabbitMQRoutingKeyKey.String("routing.key"),
					semconv.MessagingOperationReceive,
				},
				parentSpanID: trace.SpanContextFromContext(parentCtx).SpanID(),
				kind:         trace.SpanKindConsumer,
			},
		},
	}

	// start tracing all delivery spans in the tests into the span recorder
	contextPropagator := propagation.TraceContext{}
	for _, test := range tests {
		d := test.delivery

		// inject the parent context into the delivery so that the span id and trace id are
		// propagated into the delivery headers.
		injectTraceInfo(parentCtx, d)

		ctx := context.Background()
		_, endSpan := StartDeliverySpan(ctx, d,
			WithPropagators(contextPropagator),
			WithTracerProvider(provider),
		)

		endSpan()
	}

	// stop span recording.
	spans := sr.Completed()

	// assertions.
	assert.Len(t, spans, 1)

	for i, test := range tests {
		t.Run(fmt.Sprint("index", i), func(t *testing.T) {
			span := spans[i]

			assert.Equal(t, test.expected.parentSpanID, span.ParentSpanID())

			var sc trace.SpanContext
			if i == 0 {
				sc = trace.SpanContextFromContext(contextPropagator.Extract(context.Background(), NewDeliveryCarrier(test.delivery)))
			} else {
				sc = trace.SpanContextFromContext(contextPropagator.Extract(context.Background(), NewDeliveryCarrier(test.delivery)))
				sc = sc.WithRemote(false)
			}
			assert.Equal(t, sc, span.SpanContext())

			assert.Equal(t, "rabbitmq.consume", span.Name())
			assert.Equal(t, test.expected.kind, span.SpanKind())
			for _, k := range test.expected.attributeList {
				assert.Equal(t, k.Value, span.Attributes()[k.Key], k.Key)
			}
		})
	}

}

func BenchmarkStartDeliverySpan(b *testing.B) {

	delivery := amqp.Delivery{
		Body:       []byte("testing2"),
		Headers:    amqp.Table{"testkey2": "testvalue2"},
		RoutingKey: "routing.key",
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, endSpan := StartDeliverySpan(ctx, delivery)
		defer endSpan()
	}
}
