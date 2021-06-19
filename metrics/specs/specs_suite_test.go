package specs

import (
	"github.com/eyewa/eyewa-go-lib/metrics"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"testing"
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
	ml, err := metrics.NewMetricLauncher(metrics.Prometheus)
	Expect(err).Should(BeNil())
	ml.SetMeterProvider()

	ts = httptest.NewServer(http.HandlerFunc(ml.Exporter.ServeHTTP))
})

var _ = AfterSuite(func() {
	ts.Close()
})
