package tracing

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// UnaryServerTraceInterceptor intercepts unary requests on a grpc server
// and starts a new span.
func UnaryServerTraceInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor())
}

// StreamServerTraceInterceptor intercepts a stream on a grpc server
// and starts a new span.
func StreamServerTraceInterceptor() grpc.ServerOption {
	return grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor())
}
