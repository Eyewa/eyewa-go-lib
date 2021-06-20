package specs

import (
	"github.com/eyewa/eyewa-go-lib/metrics"
	"github.com/eyewa/eyewa-go-lib/metrics/prometheus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSpecs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Specs Suite")
}

var (
	ts    *httptest.Server
	errCh <-chan error
)

var _ = BeforeSuite(func() {
	option := prometheus.ExportOption{
		CollectPeriod: 1 * time.Second,
	}
	exporter, err := prometheus.NewPrometheusExporter(option)
	Expect(err).Should(BeNil())

	ml := metrics.NewMetricLauncher(exporter)
	ml.SetMeterProvider()

	ts = httptest.NewServer(http.HandlerFunc(ml.Exporter.ServeHTTP))
})

var _ = AfterSuite(func() {
	ts.Close()
})
