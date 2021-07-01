package amqp

import (
	"context"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestStartDeliverySpan(t *testing.T) {
	d := amqp.Delivery{}
	parentCtx := context.Background()
	ctx, endSpan := StartDeliverySpan(parentCtx, d)
	defer endSpan()

	assert.NotZero(t, ctx)
}
