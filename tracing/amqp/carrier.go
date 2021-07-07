package amqp

import (
	"fmt"

	"go.opentelemetry.io/otel/propagation"
)

var _ propagation.TextMapCarrier = (*HeaderCarrier)(nil)

// Get returns the value associated with the passed key.
func (hc HeaderCarrier) Get(key string) string {
	val := hc[key]
	if val != nil {
		// convert to string
		return fmt.Sprintf("%v", val)
	}
	return ""
}

// Set stores the key-value pair.
func (hc HeaderCarrier) Set(key, val string) {
	hc[key] = val
}

// Keys lists the keys stored in this carrier.
func (hc HeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range hc {
		keys = append(keys, k)
	}
	return keys
}
