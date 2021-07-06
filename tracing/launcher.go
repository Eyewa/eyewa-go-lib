package tracing

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/eyewa/eyewa-go-lib/errors"
	"github.com/eyewa/eyewa-go-lib/log"
	_ "github.com/eyewa/eyewa-go-lib/log"

	"github.com/ory/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

var (
	config          Config
	exporterTimeout = 4 * time.Second
)

// intitialises and verifies the validity of a configuration.
func initConfig() (Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// setup default configuration
	viper.SetDefault("TRACING_BLOCK_EXPORTER", "false")
	viper.SetDefault("TRACING_SECURE_EXPORTER", "false")
	hostname, err := os.Hostname()
	if err == nil && hostname != "" {
		viper.SetDefault("HOSTNAME", hostname)
	}

	envVars := []string{
		"SERVICE_NAME",
		"TRACING_BLOCK_EXPORTER",
		"TRACING_SECURE_EXPORTER",
		"TRACING_EXPORTER_ENDPOINT",
		"HOSTNAME",
	}

	for _, v := range envVars {
		if err := viper.BindEnv(v); err != nil {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if config.ServiceName == "" {
		return config, errors.ErrorNoServiceNameSpecified
	}

	if config.TracingExporterEndpoint == "" {
		return config, errors.ErrorNoExporterEndpointSpecified
	}

	log.Debug(fmt.Sprintf("Tracing config initialised: %v", config))
	return config, nil
}

// Launch launches a tracing environment and returns a
// function to shutdown.
func Launch() (ShutdownFunc, error) {
	ctx := context.Background()
	shutdownfunc := func() error {
		return nil
	}
	var err error
	_, err = initConfig()
	if err != nil {
		return shutdownfunc, fmt.Errorf("Failed to init tracing config: %v", err)
	}

	exp, err := newOtelExporter()
	if err != nil {
		return shutdownfunc, err
	}

	res, err := newResource(ctx)
	if err != nil {
		return shutdownfunc, err
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

	l := &launcher{
		exporter:  exp,
		resource:  res,
		spanprocs: processors,
		provider:  tp,
	}

	if err = l.launch(ctx); err != nil {
		log.Error(fmt.Sprintf("Failed to launch tracing launcher: %v", err))
		l.shutdown(ctx)
		return shutdownfunc, err
	}

	shutdownfunc = func() error {
		if err := l.shutdown(ctx); err != nil {
			log.Error(fmt.Sprintf("Failed to shutdown tracing launcher: %v", err))
			return err
		}
		return nil
	}

	return shutdownfunc, nil
}

// launch initiates the connection to the exporter.
func (tl *launcher) launch(ctx context.Context) error {
	log.Debug("Launching tracing launcher.")
	if err := tl.exporter.Start(ctx); err != nil {
		log.Error(fmt.Sprintf("Failed to start tracing exporter: %v", err))
		return err
	}
	return nil
}

// shutdown shuts down underlying connections.
func (tl *launcher) shutdown(ctx context.Context) error {
	var err error

	// shutdown the exporter.
	if tl.exporter != nil {
		log.Debug("Shutting down tracing exporter.")
		if err = tl.exporter.Shutdown(ctx); err != nil {
			log.Error(fmt.Sprintf("Failed to shutdown tracing exporter: %v", err))
		}
	}

	// shutdown span processors.
	if len(tl.spanprocs) > 0 {
		for _, proc := range tl.spanprocs {
			log.Debug("Shutting down tracing processors.")
			if err = proc.Shutdown(ctx); err != nil {
				log.Error(fmt.Sprintf("Failed to shutdown tracing span processor: %v", err))
			}
		}
	}

	// shutdown the tracer provider.
	if tl.provider != nil {
		log.Debug("Shutting down tracing provider")
		if err = tl.provider.Shutdown(ctx); err != nil {
			log.Error(fmt.Sprintf("Failed to shutdown tracing provider: %v", err))
		}
	}

	return err
}
