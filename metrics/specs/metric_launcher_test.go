package specs

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Given that metric launcher is launched", func() {
	Describe("When request sent to metric server without any instrumentation", func() {
		var (
			res *http.Response
			err error
		)

		It("should return http status ok", func() {
			res, err = http.Get(ts.URL)

			Expect(err).Should(BeNil())
			Expect(res.StatusCode).Should(Equal(http.StatusOK))
		})
	})
})
