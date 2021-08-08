package healthcheck

func WithEndpoint(ep string) Option {
	return func(hc *healthCheck) {
		hc.endpoint = ep
	}
}
