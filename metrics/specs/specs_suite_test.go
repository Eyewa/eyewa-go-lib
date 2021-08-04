package specs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eyewa/eyewa-go-lib/metrics"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ory/viper"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

var (
	meter *metrics.Meter
	ts    *httptest.Server
)

func TestSpecs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Specs Suite")
}

var _ = BeforeSuite(func() {
	viper.Set("METRIC_COLLECT_PERIOD", "0")
	viper.Set("SERVICE_NAME", "test-service")

	mockExporter := MockPrometheusExport()
	global.SetMeterProvider(mockExporter.MeterProvider())
	meter = metrics.NewMeter("test.meter", nil)

	ts = httptest.NewServer(http.HandlerFunc(mockExporter.ServeHTTP))
})

func MockPrometheusExport() *prometheus.Exporter {
	config := prometheus.Config{}

	c := controller.New(
		processor.New(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			export.CumulativeExportKindSelector(),
			processor.WithMemory(true),
		),
		controller.WithCollectPeriod(0),
	)

	exporter, err := prometheus.New(config, c)
	Expect(err).Should(BeNil())

	return exporter
}
