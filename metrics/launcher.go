package metrics

import (
	"github.com/eyewa/eyewa-go-lib/errors"
	"github.com/eyewa/eyewa-go-lib/log"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"net/http"
	"time"
)

// Launcher is used for serving metrics.
type Launcher struct {
	exporter                *prometheus.Exporter
	enableHostInstrument    bool
	enableRuntimeInstrument bool
}

// NewLauncher initializes MetricLauncher.
func NewLauncher(exporter *prometheus.Exporter) *Launcher {
	return &Launcher{
		exporter,
		false,
		false,
	}
}

func (ml *Launcher) SetMeterProvider() *Launcher {
	global.SetMeterProvider(ml.exporter.MeterProvider())
	return ml
}

// EnableHostInstrument enables host instrumentation
func (ml *Launcher) EnableHostInstrument() *Launcher {
	ml.enableHostInstrument = true
	return ml
}

// EnableRuntimeInstrument enables runtime instrumentation
func (ml *Launcher) EnableRuntimeInstrument() *Launcher {
	ml.enableRuntimeInstrument = true
	return ml
}

// Launch starts serving metrics. Also starts Host and Runtime instruments if they are enabled.
func (ml *Launcher) Launch() <-chan error {
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

	http.HandleFunc("/", ml.exporter.ServeHTTP)

	errCh := make(chan error)
	go func(errCh chan<- error) {
		defer close(errCh)

		errCh <- http.ListenAndServe(":2222", nil)
	}(errCh)

	return errCh
}
