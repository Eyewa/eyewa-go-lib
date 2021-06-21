package specs

import (
	"github.com/eyewa/eyewa-go-lib/metrics"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	ts    *httptest.Server
)

func TestSpecs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Specs Suite")
}

var _ = BeforeSuite(func() {
	option := metrics.ExportOption{
		CollectPeriod: 1 * time.Second,
	}
	exporter, err := metrics.NewPrometheusExporter(option)
	Expect(err).Should(BeNil())

	ml := metrics.NewMetricLauncher(exporter)
	ml.SetMeterProvider()

	ts = httptest.NewServer(http.HandlerFunc(ml.Exporter.ServeHTTP))
})

var _ = AfterSuite(func() {
	ts.Close()
})
