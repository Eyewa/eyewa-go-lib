package metrics

import (
	"github.com/eyewa/eyewa-go-lib/log"
	"github.com/eyewa/eyewa-go-lib/metrics/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"net/http"
	"time"
)

const LauncherPort = ":2222"

type ExporterType string

const (
	Prometheus ExporterType = "prometheus"
)

// Exporter is a manifest for pull based metric exporters like Prometheus
type Exporter interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	MeterProvider() metric.MeterProvider
}

// MetricLauncher is used for serving metrics. It's a wrapper for OpenTelemetry
type MetricLauncher struct {
	Exporter                Exporter
	enableHostInstrument    bool
	enableRuntimeInstrument bool
}

// NewMetricLauncher initializes OpenTelemetry Prometheus Exporter.
func NewMetricLauncher(exporterType ExporterType) (*MetricLauncher, error) {
	var (
		exporter Exporter
		err      error
	)

	switch string(exporterType) {
	case string(Prometheus):
		option := prometheus.ExportOption{
			CollectPeriod: 1 * time.Second,
		}

		exporter, err = prometheus.NewPrometheusExporter(option)
		if err != nil {
			return nil, err
		}
	}

	return &MetricLauncher{
		exporter,
		false,
		false,
	}, nil
}

// SetMeterProvider sets prometheus meter provider globally
func (ml *MetricLauncher) SetMeterProvider() *MetricLauncher {
	global.SetMeterProvider(ml.Exporter.MeterProvider())
	return ml
}

// EnableHostInstrument enables host instrumentation
func (ml *MetricLauncher) EnableHostInstrument() *MetricLauncher {
	ml.enableHostInstrument = true
	return ml
}

// EnableRuntimeInstrument enables runtime instrumentation
func (ml *MetricLauncher) EnableRuntimeInstrument() *MetricLauncher {
	ml.enableRuntimeInstrument = true
	return ml
}

// Launch starts serving metrics. Also starts Host and Runtime instruments if they are enabled.
func (ml MetricLauncher) Launch() <-chan error {
	if ml.enableHostInstrument {
		err := host.Start()
		if err != nil {
			log.Fatal(FailedToStartRuntimeMetricsError.Inner(err).Error())
		}
	}

	if ml.enableRuntimeInstrument {
		err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
		if err != nil {
			log.Fatal(FailedToStartHostMetricsError.Inner(err).Error())
		}
	}

	http.HandleFunc("/", ml.Exporter.ServeHTTP)

	errCh := make(chan error)
	go func(errCh chan<- error) {
		defer close(errCh)

		errCh <- http.ListenAndServe(LauncherPort, nil)
	}(errCh)

	return errCh
}
