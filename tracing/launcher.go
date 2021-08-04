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
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
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

	log.Debug(fmt.Sprintf("Tracing config initialised: %+v \n", config))
	return config, nil
}

// Launch launches a tracing environment that that enables
// tracing on all instrumented golib packages and returns a
// function to shutdown.
func Launch() (ShutdownFunc, error) {
	ctx := context.Background()
	shutdownfunc := func() error {
		return nil
	}

	// initialize the global configuration
	_, err := initConfig()
	if err != nil {
		return shutdownfunc, fmt.Errorf("Failed to initialize tracing configuration: %v", err)
	}

	// setup and connect to the open telemetry collector
	exp, err := newOtelExporter(ctx)
	shutdownfunc = func() error {
		if err := exp.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	}

	if err != nil {
		return shutdownfunc, err
	}

	// setup the resource from which the trace comes from
	res, err := newResource(ctx)
	if err != nil {
		return shutdownfunc, err
	}

	// setup a tracer provider
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exp)),
	)

	// set globals for propagation and the trace provider
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tracerProvider)

	shutdownfunc = func() error {
		var err error

		// shutdown the tracer provider.
		// this already shuts down all underlying processors.
		if tracerProvider != nil {
			log.Debug("Shutting down tracing provider")
			if err = tracerProvider.Shutdown(ctx); err != nil {
				log.Error(fmt.Sprintf("Failed to shutdown tracing provider: %v", err))
			}
		}

		// shutdown the exporter.
		if exp != nil {
			log.Debug("Shutting down tracing exporter.")
			if err = exp.Shutdown(ctx); err != nil {
				log.Error(fmt.Sprintf("Failed to shutdown tracing exporter: %v", err))
			}
		}

		return err
	}

	return shutdownfunc, nil
}
