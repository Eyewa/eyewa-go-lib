package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/propagation"
)

var _ propagation.TextMapCarrier = (*PublishingCarrier)(nil)
var _ propagation.TextMapCarrier = (*DeliveryCarrier)(nil)

// DeliveryCarrier injects and extracts
// traces from the headers of a amqp.Delivery.
type DeliveryCarrier struct {
	delivery *amqp.Delivery
}

// NewDeliveryCarrier constructs a new DeliveryCarrier.
func NewDeliveryCarrier(d *amqp.Delivery) DeliveryCarrier {
	return DeliveryCarrier{delivery: d}
}

// Get gets a header from the delivery.
func (c *DeliveryCarrier) Get(key string) string {
	for k, h := range c.delivery.Headers {
		if h != "" && k == key {
			return fmt.Sprintf("%v", h)
		}
	}
	return ""
}

// Set sets a header on the delivery.
func (c *DeliveryCarrier) Set(key, val string) {
	for k, _ := range c.delivery.Headers {
		c.delivery.Headers[k] = val
	}
}

// Keys returns all the header keys of the delivery.
func (c *DeliveryCarrier) Keys() []string {
	keys := make([]string, len(c.delivery.Headers))
	for k, _ := range c.delivery.Headers {
		keys = append(keys, k)
	}
	return keys
}

// PublishingCarrier injects and extracts
// traces from the headers of a amqp.Publishing.
type PublishingCarrier struct {
	publishing *amqp.Publishing
}

// NewPublishingCarrier constructs a new PublishingCarrier.
func NewPublishingCarrier(p *amqp.Publishing) PublishingCarrier {
	return PublishingCarrier{publishing: p}
}

// Get gets a header from the publishing.
func (c *PublishingCarrier) Get(key string) string {
	for k, h := range c.publishing.Headers {
		if h != "" && k == key {
			return fmt.Sprintf("%v", h)
		}
	}
	return ""
}

// Set sets a header on the publishing.
func (c *PublishingCarrier) Set(key, val string) {
	for k := range c.publishing.Headers {
		c.publishing.Headers[k] = val
	}
}

// Keys returns all the header keys of the publishing.
func (c *PublishingCarrier) Keys() []string {
	keys := make([]string, len(c.publishing.Headers))
	for k := range c.publishing.Headers {
		keys = append(keys, k)
	}
	return keys
}
