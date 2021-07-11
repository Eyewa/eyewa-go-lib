package metrics

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/eyewa/eyewa-go-lib/errors"
	"github.com/eyewa/eyewa-go-lib/log"
	"github.com/ory/viper"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/metric/global"
)

func init() {
	l, err := newLauncher()
	if err != nil {
		log.Error(fmt.Sprintf(errors.ErrorFailedToStartMetricServer.Error(), err.Error()))

		return
	}

	l.setMeterProvider().
		enableHostInstrumentation().
		enableRuntimeInstrumentation().
		launch()
}

// newLauncher initializes launcher.
func newLauncher() (*launcher, error) {
	option, err := initConfig()
	if err != nil {
		return nil, err
	}

	exporter, err := newPrometheusExporter(option)
	if err != nil {
		return nil, err
	}

	return &launcher{
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

	log.SetLogLevel()

	envVars := []string{
		"METRICS_COLLECTOR_INTERVAL",
		"SERVICE_NAME",
	}

	for _, v := range envVars {
		if err := viper.BindEnv(v); err != nil {
			return exportOption, err
		}
	}

	err := viper.Unmarshal(&exportOption)
	if err != nil {
		return exportOption, err
	}

	return exportOption, nil
}

func (l *launcher) setMeterProvider() *launcher {
	global.SetMeterProvider(l.exporter.MeterProvider())

	return l
}

// enableHostInstrumentation enables host instrumentation
func (l *launcher) enableHostInstrumentation() *launcher {
	l.enableHostInstrument = true

	return l
}

// enableRuntimeInstrumentation enables runtime instrumentation
func (l *launcher) enableRuntimeInstrumentation() *launcher {
	l.enableRuntimeInstrument = true

	return l
}

// launch starts serving metrics. Also starts Host and Runtime instruments if they are enabled.
func (l *launcher) launch() {
	if l.enableHostInstrument {
		err := host.Start()
		if err != nil {
			log.Error(errors.ErrorFailedToStartRuntimeMetrics.Error())
		}
	}

	if l.enableRuntimeInstrument {
		err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
		if err != nil {
			log.Error(errors.ErrorFailedToStartHostMetrics.Error())
		}
	}

	http.HandleFunc("/", l.exporter.ServeHTTP)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error(fmt.Sprintf(errors.ErrorFailedToStartMetricServer.Error(), r.(error).Error()))
			}
		}()

		err := http.ListenAndServe(":2222", nil)
		if err != nil {
			log.Error(fmt.Sprintf(errors.ErrorFailedToStartMetricServer.Error(), err.Error()))
		}
	}()
}
