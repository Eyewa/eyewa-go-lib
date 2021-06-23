package trace

// "go.opentelemetry.io/otel/oteltest"

// func TestResource(t *testing.T) {

// 	assert.
// }
// func TestInterceptors(t *testing.T) {
// 	clientUnarySR := new(oteltest.SpanRecorder)
// 	clientUnaryTP := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(clientUnarySR))

// 	clientStreamSR := new(oteltest.SpanRecorder)
// 	clientStreamTP := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(clientStreamSR))

// 	serverUnarySR := new(oteltest.SpanRecorder)
// 	serverUnaryTP := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(serverUnarySR))

// 	serverStreamSR := new(oteltest.SpanRecorder)
// 	serverStreamTP := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(serverStreamSR))

// 	assert.NoError(t, doCalls(
// 		[]grpc.DialOption{
// 			grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor(otelgrpc.WithTracerProvider(clientUnaryTP))),
// 			grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor(otelgrpc.WithTracerProvider(clientStreamTP))),
// 		},
// 		[]grpc.ServerOption{
// 			grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(serverUnaryTP))),
// 			grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor(otelgrpc.WithTracerProvider(serverStreamTP))),
// 		},
// 	))

// 	t.Run("UnaryClientSpans", func(t *testing.T) {
// 		checkUnaryClientSpans(t, clientUnarySR.Completed())
// 	})

// 	t.Run("StreamClientSpans", func(t *testing.T) {
// 		checkStreamClientSpans(t, clientStreamSR.Completed())
// 	})

// 	t.Run("UnaryServerSpans", func(t *testing.T) {
// 		checkUnaryServerSpans(t, serverUnarySR.Completed())
// 	})

// 	t.Run("StreamServerSpans", func(t *testing.T) {
// 		checkStreamServerSpans(t, serverStreamSR.Completed())
// 	})
// }
