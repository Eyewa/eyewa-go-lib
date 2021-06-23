package specs

import (
	"github.com/eyewa/eyewa-go-lib/metrics"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	ts    *httptest.Server
	meter *metrics.Meter
)

func TestSpecs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Specs Suite")
}

var _ = BeforeSuite(func() {
	option := metrics.ExportOption{
		CollectPeriod: 0,
	}
	exporter, err := metrics.NewPrometheusExporter(option)
	Expect(err).Should(BeNil())

	ml := metrics.NewLauncher(exporter)
	ml.SetMeterProvider()

	meter = metrics.NewMeter("test.meter", nil)

	ts = httptest.NewServer(http.HandlerFunc(exporter.ServeHTTP))
})

var _ = AfterSuite(func() {
	ts.Close()
})
