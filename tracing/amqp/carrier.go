package amqp

import (
	"fmt"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/propagation"
)

var _ propagation.TextMapCarrier = (*PublishingCarrier)(nil)
var _ propagation.TextMapCarrier = (*DeliveryCarrier)(nil)

// DeliveryCarrier constructs a new DeliveryCarrier.
func NewDeliveryCarrier(d amqp.Delivery) DeliveryCarrier {
	return DeliveryCarrier{delivery: d}
}

// Get gets a header from the delivery.
func (c DeliveryCarrier) Get(key string) string {
	val := c.delivery.Headers[key]
	if val != nil {
		// convert to string
		return fmt.Sprintf("%v", val)
	}

	return ""
}

// Set sets a header on the delivery.
func (c DeliveryCarrier) Set(key, val string) {
	c.delivery.Headers[key] = val
}

// Keys returns all the header keys of the delivery.
func (c DeliveryCarrier) Keys() []string {
	var keys []string
	for k := range c.delivery.Headers {
		keys = append(keys, k)
	}

	if len(keys) > 0 {
		return keys
	}

	return []string{}
}

// NewPublishingCarrier constructs a new PublishingCarrier.
func NewPublishingCarrier(p amqp.Publishing) PublishingCarrier {
	return PublishingCarrier{publishing: p}
}

// Get gets a header from the publishing.
func (c PublishingCarrier) Get(key string) string {
	val := c.publishing.Headers[key]
	if val != nil {
		// convert to string
		return fmt.Sprintf("%v", val)
	}

	return ""
}

// Set sets a header on the publishing.
func (c PublishingCarrier) Set(key, val string) {
	c.publishing.Headers[key] = val
}

// Keys returns all the header keys of the publishing.
func (c PublishingCarrier) Keys() []string {
	var keys []string
	for k := range c.publishing.Headers {
		keys = append(keys, k)
	}

	if len(keys) > 0 {
		return keys
	}

	return []string{}
}
