// +build !mock

package tracing

import (
	"context"

	libErrs "github.com/eyewa/eyewa-go-lib/errors"
	exporter "github.com/eyewa/eyewa-go-lib/tracing/exporter"
)

// func initConfig() (Config, error) {
// 	viper.AutomaticEnv()
// 	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

// 	envVars := []string{
// 		"SERVICE_NAME",
// 		"HOST_NAME",
// 		"TRACE_COLLECTOR_ENDPOINT",
// 	}

// 	for _, v := range envVars {
// 		if err := viper.BindEnv(v); err != nil {
// 			return config, err
// 		}
// 	}
// 	if err := viper.Unmarshal(&config); err != nil {
// 		return config, err
// 	}

// 	var c Config
// 	envError := envconfig.Process(context.Background(), &c)
// 	c.logger = &DefaultLogger{}
// 	c.errorHandler = &defaultHandler{logger: c.logger}
// 	var defaultOpts []Option

// 	for _, opt := range append(defaultOpts, opts...) {
// 		opt(&c)
// 	}
// 	c.Resource = newResource(&c)

// 	if envError != nil {
// 		c.logger.Fatalf("environment error: %v", envError)
// 	}

// 	return c

// 	return config, nil
// }

// // Configure configures the environment for tracing. It inits env vars first
// // and injected options as overrides then returns a shutdown function.
// func Configure(opts ...Option) Connector {
// 	cfg := Config{}
// 	connector := Connector{cfg}
// 	connector.Connect()
// 	c := newConfig(opts...)

// 	err := validateConfiguration(c)
// 	if err != nil {
// 		log.Error(fmt.Sprintf("configuration error: %v", err))
// 	}

// 	l := Connector{
// 		config: c,
// 	}
// 	for _, setup := range []setupFunc{setupTracing, setupMetrics} {
// 		shutdown, err := setup(c)
// 		if err != nil {
// 			c.logger.Fatalf("setup error: %v", err)
// 			continue
// 		}
// 		if shutdown != nil {
// 			l.shutdownFuncs = append(l.shutdownFuncs, shutdown)
// 		}
// 	}
// 	return l
// }

// // ConfigureAndConnect configures and connects to the tracing backend using
// // the trace configuration. It returns a graceful shutdown function.
// func ConfigureAndConnect() (func(), error) {

// 	ctx := context.Background()

// 	// initialize env var configs
// 	_, err := initConfig()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// configure the exporter
// 	ep := config.CollectorEndpoint
// 	var exp *otlp.Exporter
// 	if config.CollectorEndpoint != "" {
// 		driver := newDriver(ep)
// 		exp = newExporter(driver)
// 		exp.Start()
// 		exporter, err = stdout.NewExporter(stdout.WithPrettyPrint())
// 	} else {

// 	}

// 	if err != nil { // TODO: wrap the error
// 		return nil, libErr.ErrorFailedToCreateTraceExporter
// 	}

// 	// configure the resource with semantic conversions
// 	host := config.HostName
// 	svcname := config.ServiceName
// 	opts := resource.WithAttributes(
// 		semconv.HostNameKey.String(host),
// 		semconv.ServiceNameKey.String(svcname),
// 	)
// 	res, err := resource.New(ctx, opts)
// 	if err != nil { // TODO: wrap the error
// 		return nil, errors.New(fmt.Sprintf("couldnt create a new trace resource for endpoint: %s", ep))
// 	}

// 	tp := newProvider(exporter, res)
// 	registerTraceProvider(tp)

// 	return func() {
// 		err := tp.Shutdown(ctx)
// 		if err != nil {
// 			log.Error("failed to shutdown trace provider")
// 		}
// 		err = exporter.Shutdown(ctx)
// 		if err != nil {
// 			log.Error("failed to shutdown trace exporter")
// 		}
// 	}, nil
// }

// // registerTraceProvider registers the trace provider globally
// func registerTraceProvider(tp *sdktrace.TracerProvider) {
// 	otel.SetTextMapPropagator(propagation.TraceContext{})
// 	otel.SetTracerProvider(tp)
// }

// // newProvider creates a new trace provider with default options for the given resource and exporter.
// func newProvider(exp trace.SpanExporter, res *resource.Resource) *sdktrace.TracerProvider {
// 	bsp := sdktrace.NewBatchSpanProcessor(exp)
// 	samp := sdktrace.AlwaysSample()
// 	return sdktrace.NewTracerProvider(
// 		sdktrace.WithSampler(samp),
// 		sdktrace.WithResource(res),
// 		sdktrace.WithSpanProcessor(bsp),
// 	)
// }

// New launcher constructs a new launcher that sets
// up a environment for tracing.
func NewLauncher(e exporter.Exporter) *Launcher {
	return &Launcher{exporter: e, ctx: context.Background()}
}

// Launch launches the launcher and connects to exporter.
func (l *Launcher) Launch() error {
	if err := l.exporter.Start(l.ctx); err != nil {
		return libErrs.Wrap(err, libErrs.ErrExporterStartupFailure)
	}

	return nil
}

// Shutdown shuts down the launcher and all connections.
func (l *Launcher) Shutdown() error {
	if err := l.exporter.Shutdown(l.ctx); err != nil {
		return err
	}

	return nil
}

// Launches a tracing environment by constructing a new
// stdout exporter when no exporter selected. Setting up the host
// resource and setting the global trace provider.
func Launch() (ShutdownFunc, error) {
	exp, err := exporter.NewStdOut()
	if err != nil {
		return nil, err
	}
	launcher := NewLauncher(exp)
	if err := launcher.Launch(); err != nil {
		return nil, err
	}

	return launcher.Shutdown, nil
}
