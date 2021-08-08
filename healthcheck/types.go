package healthcheck

import "github.com/InVisionApp/go-health"

// Option configures a HealthCheck.
type Option func(s *healthCheck)

// healthCheck wraps go-health health instance.
type healthCheck struct {
	*health.Health
	endpoint string
	addr     string
}
