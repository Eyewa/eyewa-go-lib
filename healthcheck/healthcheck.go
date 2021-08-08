package healthcheck

import (
	"net/http"

	"github.com/InVisionApp/go-health"
	"github.com/InVisionApp/go-health/handlers"
)

// Start starts checking the health of dependencies
// and configures a default http endpoint that returns
// the health status of those dependencies.
func Start(opts ...Option) error {
	hc := new(opts...)

	// start executing asynchronous checks
	if err := hc.start(); err != nil {
		return err
	}

	http.HandleFunc(hc.endpoint, handlers.NewJSONHandlerFunc(hc, nil))
	return http.ListenAndServe(hc.addr, nil)
}

// New constructs a new HealthCheck instance.
func new(opts ...Option) *healthCheck {
	const (
		defaultEndpoint = "/healthcheck"
		defaultAddr     = ":3333"
	)

	hc := &healthCheck{
		health.New(),
		defaultEndpoint,
		defaultAddr,
	}

	// TODO: get the mysql + postgres db connections to allow pinging.
	// mysqlCheck, err := checkers.NewSQL(&checkers.SQLConfig{
	// 	Pinger:
	// })

	// if err != nil {
	// 	return err
	// }

	// Add the default checks
	// hc.AddChecks([]*health.Config{
	// 	{
	// 		Name:     "sql-check",
	// 		Checker:  sqlCheck,
	// 		Interval: time.Duration(3) * time.Second,
	// 		Fatal:    true,
	// 	},
	// })

	for _, opt := range opts {
		opt(hc)
	}

	return hc
}

// Start starts checking the health of configured dependencies.
func (hc *healthCheck) start() error {
	return hc.Start()
}
