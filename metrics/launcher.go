package metrics

import (
	"github.com/eyewa/eyewa-go-lib/errors"
	"github.com/eyewa/eyewa-go-lib/log"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/metric/global"
	"net/http"
	"time"
)

// MetricLauncher is used for serving metrics.
type MetricLauncher struct {
	Exporter                *PrometheusExporter
	enableHostInstrument    bool
	enableRuntimeInstrument bool
}

// NewMetricLauncher initializes MetricLauncher.
func NewMetricLauncher(exporter *PrometheusExporter) *MetricLauncher {
	return &MetricLauncher{
		exporter,
		false,
		false,
	}
}

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
			log.Error(errors.FailedToStartRuntimeMetricsError.Error())
		}
	}

	if ml.enableRuntimeInstrument {
		err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
		if err != nil {
			log.Error(errors.FailedToStartHostMetricsError.Error())
		}
	}

	http.HandleFunc("/", ml.Exporter.ServeHTTP)

	errCh := make(chan error)
	go func(errCh chan<- error) {
		defer close(errCh)

		errCh <- http.ListenAndServe(":2222", nil)
	}(errCh)

	return errCh
}
