package specs

import (
	"github.com/eyewa/eyewa-go-lib/metrics"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ory/viper"
	"testing"
)

var (
	meter *metrics.Meter
	URL   = "http://localhost:2222"
)

func TestSpecs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Specs Suite")
}

var _ = BeforeSuite(func() {
	viper.Set("METRIC_COLLECT_PERIOD", "0")
	viper.Set("SERVICE_NAME", "test-service")

	ml, err := metrics.NewLauncher()
	Expect(err).Should(BeNil())

	ml.SetMeterProvider().Launch()

	meter = metrics.NewMeter("test.meter", nil)
})
