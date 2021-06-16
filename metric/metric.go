package metric

import (
	"errors"
	"fmt"
	"github.com/eyewa/eyewa-go-lib/log"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"net/http"
	"time"
)

const Port = ":2222"

// startHostInstrument starts Host instrumentation.
// https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/host@v0.20.0
func startHostInstrument() error {
	err := host.Start()
	if err != nil {
		return err
	}
	return nil
}

// startHostInstrument starts Runtime instrumentation.
// https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/runtime@v0.20.0
func startRuntimeInstrument() error {
	err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		return err
	}
	return nil
}

// Option is for metric package initiation.
// CollectPeriod sets period interval exporter.
type Option struct {
	CollectPeriod time.Duration
}

// Launch initializes OpenTelemetry Prometheus Exporter.
// Also starts Host and Runtime instruments.
func Launch(option Option) error {
	config := prometheus.Config{}
	c := controller.New(
		processor.New(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			export.CumulativeExportKindSelector(),
			processor.WithMemory(true),
		),
		controller.WithCollectPeriod(option.CollectPeriod),
	)

	exporter, err := prometheus.New(config, c)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to initialize prometheus exporter %v", err))
	}
	global.SetMeterProvider(exporter.MeterProvider())

	if err := startHostInstrument(); err != nil {
		log.Fatal(fmt.Sprintf("failed to start runtime metrics: %v", err))
	}

	if err := startRuntimeInstrument(); err != nil {
		log.Fatal(fmt.Sprintf("failed to start host metrics: %v", err))
	}

	http.HandleFunc("/", exporter.ServeHTTP)
	go func() {
		err := http.ListenAndServe(Port, nil)
		if err != nil {
			log.Fatal(fmt.Sprintf("failed to start metric server: %v", err))
		}
	}()

	return nil
}
