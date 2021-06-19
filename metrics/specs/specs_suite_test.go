package specs
import (
	"github.com/eyewa/eyewa-go-lib/metrics"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSpecs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Specs Suite")
}

var (
	ts *httptest.Server
	errCh <-chan error
)

var _ = BeforeSuite(func() {
	option := metrics.ExportOption{
		CollectPeriod: 10 * time.Millisecond,
	}

	ml, err := metrics.NewMetricLauncher(option)
	Expect(err).Should(BeNil())
	ml.SetMeterProvider()

	ts = httptest.NewServer(http.HandlerFunc(ml.Exporter.ServeHTTP))
})

var _ = AfterSuite(func() {
	ts.Close()
})
