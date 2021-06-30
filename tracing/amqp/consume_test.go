package amqp

import (
	"errors"
	"fmt"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestWrapConsume(t *testing.T) {
	// prepare channel with fake delivery
	fakeDelivery := amqp.Delivery{RoutingKey: "test.routing.key"}
	out := make(chan amqp.Delivery, 1)
	out <- fakeDelivery
	consume := func(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
		fmt.Print("Called consume")
		return out, nil
	}
	close(out)

	wrappedConsume := WrapConsume(consume)
	msgs, err := wrappedConsume("testqueue", "testconsumer", false, false, false, false, nil)

	for d := range msgs {
		assert.Equal(t, d, fakeDelivery)
		break
	}

	assert.Nil(t, err)
	assert.NoError(t, err)
	assert.Zero(t, err)
}

func TestWrapConsumeFail(t *testing.T) {
	consume := func(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
		var out chan amqp.Delivery // channel zero value
		return out, errors.New("failure")
	}

	wrappedConsume := WrapConsume(consume)
	msgs, err := wrappedConsume("testqueue", "testconsumer", false, false, false, false, nil)

	assert.NotNil(t, err)
	assert.Error(t, err)

	assert.Zero(t, msgs)
}
