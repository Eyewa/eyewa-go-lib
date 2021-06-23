package trace

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/eyewa/eyewa-go-lib/log"
	"github.com/ory/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	stdout "go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"google.golang.org/grpc"
)

var (
	config   EnvConfig
	exporter otlp.Exporter
)

func initConfig() (EnvConfig, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	envVars := []string{
		"SERVICE_NAME",
		"HOST_NAME",
		"TRACE_EXPORTER_ENDPOINT",
	}

	for _, v := range envVars {
		if err := viper.BindEnv(v); err != nil {
			return config, err
		}
	}
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}

// ConfigureAndConnect configures and connects to the tracing backend using
// the trace configuration. It returns a graceful shutdown function.
func ConfigureAndConnect() (func(), error) {
	ctx := context.Background()

	// initialize env var configs
	_, err := initConfig()
	if err != nil {
		return nil, err
	}

	// configure the exporter
	ep := config.CollectorEndpoint
	exp, err := configureExporter(ctx, ep)
	if err != nil { // TODO: wrap the error
		return nil, errors.New(fmt.Sprintf("couldnt create a trace exporter for endpoint %s", ep))
	}

	// configure the resource with semantic conversions
	host := config.HostName
	svcname := config.ServiceName
	opts := resource.WithAttributes(
		semconv.HostNameKey.String(host),
		semconv.ServiceNameKey.String(svcname),
	)
	res, err := resource.New(ctx, opts)
	if err != nil { // TODO: wrap the error
		return nil, errors.New(fmt.Sprintf("couldnt create a new trace resource for endpoint: %s", ep))
	}

	tp := newProvider(exp, res)
	registerTraceProvider(tp)

	return func() {
		err := tp.Shutdown(ctx)
		if err != nil {
			log.Error("failed to shutdown trace provider")
		}
		err = exp.Shutdown(ctx)
		if err != nil {
			log.Error("failed to shutdown trace exporter")
		}
	}, nil
}

// registerTraceProvider registers the trace provider globally
func registerTraceProvider(tp *sdktrace.TracerProvider) {
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)
}

// instantiates a new collector exporter and defaults to stdout if no enpoint is provided
// returns the connect function
func configureExporter(ctx context.Context, endpoint string) (sdktrace.SpanExporter, error) {
	if endpoint == "" {
		r, err := stdout.NewExporter(stdout.WithPrettyPrint())
		return r, err
	}

	driver := otlpgrpc.NewDriver(
		// need to look at securing this connection later
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(endpoint),
		otlpgrpc.WithDialOption(grpc.WithBlock()), // useful for testing
	)

	// create a new exporter that connects to the endpoint
	exp, err := otlp.NewExporter(ctx, driver)
	return exp, err
}

// newProvider creates a new trace provider for the given resource and exporter.
func newProvider(exp trace.SpanExporter, res *resource.Resource) *sdktrace.TracerProvider {
	bsp := sdktrace.NewBatchSpanProcessor(exp)
	samp := sdktrace.AlwaysSample()
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(samp),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	return tracerProvider
}
