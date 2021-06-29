package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/propagation"
)

var _ propagation.TextMapCarrier = (*PublishingHeaderCarrier)(nil)
var _ propagation.TextMapCarrier = (*DeliveryHeaderCarrier)(nil)

// DeliveryHeaderCarrier injects and extracts
// traces from the headers of a amqp.Delivery.
type DeliveryHeaderCarrier struct {
	delivery *amqp.Delivery
}

// NewDeliveryHeaderCarrier constructs a new DeliveryHeaderCarrier.
func NewDeliveryHeaderCarrier(d *amqp.Delivery) DeliveryHeaderCarrier {
	return DeliveryHeaderCarrier{delivery: d}
}

// Get gets a header from the delivery.
func (c *DeliveryHeaderCarrier) Get(key string) string {
	val := c.delivery.Headers[key]
	if val != nil {
		// convert to string
		return fmt.Sprintf("%v", val)
	}

	return ""
}

// Set sets a header on the delivery.
func (c *DeliveryHeaderCarrier) Set(key, val string) {
	c.delivery.Headers[key] = val
}

// Keys returns all the header keys of the delivery.
func (c *DeliveryHeaderCarrier) Keys() []string {
	var keys []string
	for k := range c.delivery.Headers {
		keys = append(keys, k)
	}

	if len(keys) > 0 {
		return keys
	}

	return []string{}
}

// PublishingHeaderCarrier injects and extracts
// traces from the headers of a amqp.Publishing.
type PublishingHeaderCarrier struct {
	publishing *amqp.Publishing
}

// NewPublishingHeaderCarrier constructs a new PublishingHeaderCarrier.
func NewPublishingHeaderCarrier(p *amqp.Publishing) PublishingHeaderCarrier {
	return PublishingHeaderCarrier{publishing: p}
}

// Get gets a header from the publishing.
func (c *PublishingHeaderCarrier) Get(key string) string {
	val := c.publishing.Headers[key]
	if val != nil {
		// convert to string
		return fmt.Sprintf("%v", val)
	}

	return ""
}

// Set sets a header on the publishing.
func (c *PublishingHeaderCarrier) Set(key, val string) {
	c.publishing.Headers[key] = val
}

// Keys returns all the header keys of the publishing.
func (c *PublishingHeaderCarrier) Keys() []string {
	var keys []string
	for k := range c.publishing.Headers {
		keys = append(keys, k)
	}

	if len(keys) > 0 {
		return keys
	}

	return []string{}
}
