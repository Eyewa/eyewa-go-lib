package specs

import (
	"github.com/eyewa/eyewa-go-lib/metrics"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
	option := metrics.ExportOption{
		CollectPeriod: 0,
	}
	ml, err := metrics.NewLauncher(option)
	Expect(err).Should(BeNil())

	ml.SetMeterProvider().Launch()

	meter = metrics.NewMeter("test.meter", nil)
})
