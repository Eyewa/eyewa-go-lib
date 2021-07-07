package amqp

// import (
// 	"go.opentelemetry.io/contrib"
// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/propagation"
// 	"go.opentelemetry.io/otel/trace"
// )

// var (
// 	instrumentationName = "github.com/eyewa/eyewa-go-lib/tracing/amqp"
// 	messagingSystem     = "rabbitmq"
// )

// // newConfig returns a config with all Options set.
// func newConfig(opts ...Option) config {
// 	cfg := config{
// 		Propagators:    otel.GetTextMapPropagator(),
// 		TracerProvider: otel.GetTracerProvider(),
// 	}
// 	for _, opt := range opts {
// 		opt(&cfg)
// 	}

// 	cfg.Tracer = cfg.TracerProvider.Tracer(
// 		instrumentationName,
// 		trace.WithInstrumentationVersion(contrib.SemVersion()),
// 	)

// 	return cfg
// }

// // WithTracerProvider specifies a tracer provider to use for creating a tracer.
// // If none is specified, the global provider is used.
// func WithTracerProvider(provider trace.TracerProvider) Option {
// 	return func(cfg *config) {
// 		cfg.TracerProvider = provider
// 	}
// }

// // WithPropagators specifies propagators to use for extracting
// // information from the amqp messages. If none are specified, global
// // ones will be used.
// func WithPropagators(propagators propagation.TextMapPropagator) Option {
// 	return func(cfg *config) {
// 		cfg.Propagators = propagators
// 	}
// }
