package tracing

import (
	"context"
	"fmt"
	"strings"

	"github.com/ory/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

var (
	l *launcher
)

func initConfig() (config, error) {
	var config config

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("EXPORTER_BLOCKING", "false")
	viper.SetDefault("EXPORTER_SECURE", "true")

	envVars := []string{
		"EXPORTER_BLOCKING",
		"EXPORTER_SECURE",
		"EXPORTER_ENDPOINT",
		"SERVICE_VERSION",
		"SERVICE_NAME",
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

// Launch launches a tracing environment and returns a
// function to shutdown.
func Launch() (ShutdownFunc, error) {
	config, err := initConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed to init config: %v", err)
	}

	exp := newOtelCollectorExporter(
		config.ExporterEndpoint,
		config.ExporterSecure,
		config.ExporterBlocking,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create span exporter: %v", err)
	}

	res, err := newResource(
		config.ServiceName,
		config.ServiceVersion,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create a resource: %v", err)
	}

	var processors []trace.SpanProcessor
	bsp := trace.NewBatchSpanProcessor(exp)
	processors = append(processors, bsp)

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)

	registerPropagators()

	otel.SetTracerProvider(tp)

	l = &launcher{
		config:    config,
		exporter:  exp,
		resource:  res,
		spanprocs: processors,
		provider:  tp,
	}

	ctx := context.Background()
	if err = l.launch(ctx); err != nil {
		return nil, err
	}

	shutdownfunc := func() error {
		if err := l.shutdown(ctx); err != nil {
			return err
		}
		return nil
	}

	return shutdownfunc, nil
}

// launch initiates the connection to the exporter.
func (l *launcher) launch(ctx context.Context) error {
	if err := l.exporter.Start(ctx); err != nil {
		return err
	}
	return nil
}

// shutdown shuts down underlying connections.
func (l *launcher) shutdown(ctx context.Context) error {
	var err error
	err = l.exporter.Shutdown(ctx)
	for _, proc := range l.spanprocs {
		err = proc.Shutdown(ctx)
	}
	return err
}
