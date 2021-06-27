package metrics

import (
	"github.com/eyewa/eyewa-go-lib/errors"
	"github.com/eyewa/eyewa-go-lib/log"
	"github.com/ory/viper"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/metric/global"
	"net/http"
	"strings"
	"time"
)

// NewLauncher initializes MetricLauncher.
func NewLauncher() (*Launcher, error) {
	option, err := initConfig()
	if err != nil {
		return nil, err
	}

	exporter, err := newPrometheusExporter(option)
	if err != nil {
		return nil, err
	}

	return &Launcher{
		exporter,
		false,
		false,
	}, nil
}

func initConfig() (ExportOption, error) {
	var exportOption ExportOption

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("METRICS_COLLECTOR_INTERVAL", "10s")

	envVars := []string{
		"METRICS_COLLECTOR_INTERVAL",
	}

	for _, v := range envVars {
		if err := viper.BindEnv(v); err != nil {
			return exportOption, err
		}
	}
	if err := viper.Unmarshal(&exportOption); err != nil {
		return exportOption, err
	}

	return exportOption, nil
}

func (ml *Launcher) SetMeterProvider() *Launcher {
	global.SetMeterProvider(ml.exporter.MeterProvider())
	return ml
}

// EnableHostInstrumentation enables host instrumentation
func (ml *Launcher) EnableHostInstrumentation() *Launcher {
	ml.enableHostInstrument = true
	return ml
}

// EnableRuntimeInstrumentation enables runtime instrumentation
func (ml *Launcher) EnableRuntimeInstrumentation() *Launcher {
	ml.enableRuntimeInstrument = true
	return ml
}

// Launch starts serving metrics. Also starts Host and Runtime instruments if they are enabled.
func (ml *Launcher) Launch() {
	if ml.enableHostInstrument {
		err := host.Start()
		if err != nil {
			log.Error(errors.ErrorFailedToStartRuntimeMetrics.Error())
		}
	}

	if ml.enableRuntimeInstrument {
		err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
		if err != nil {
			log.Error(errors.ErrorFailedToStartHostMetrics.Error())
		}
	}

	http.HandleFunc("/", ml.exporter.ServeHTTP)

	go func() {
		err := http.ListenAndServe(":2222", nil)
		if err != nil {
			log.Error(errors.ErrorFailedToStartMetricServer.Error())
		}
	}()
}
