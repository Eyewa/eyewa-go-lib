package tracing

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// UnaryClientTraceInterceptor intercepts requests from a client to a grpc server
// and starts a span
func UnaryClientTraceInterceptor() grpc.DialOption {
	return grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor())
}

// StreamClientTraceInterceptor intercepts a stream from client to a grpc server
// and starts a span
func StreamClientTraceInterceptor() grpc.DialOption {
	return grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor())
}
