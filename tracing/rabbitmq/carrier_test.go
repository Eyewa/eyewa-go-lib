package rabbitmq

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestPublishingCarrierGet(t *testing.T) {
	testCases := []struct {
		name     string
		carrier  PublishingCarrier
		key      string
		expected string
	}{
		{
			name: "exists",
			carrier: PublishingCarrier{publishing: &amqp.Publishing{Headers: amqp.Table{
				"foo": "bar",
			}}},
			key:      "foo",
			expected: "bar",
		},
		{
			name:     "not exists",
			carrier:  PublishingCarrier{publishing: &amqp.Publishing{Headers: amqp.Table{}}},
			key:      "foo",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.carrier.Get(tc.key)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestPublishingCarrierSet(t *testing.T) {
	pub := amqp.Publishing{Headers: amqp.Table{
		"foo": "bar",
	}}
	carrier := PublishingCarrier{publishing: &pub}

	carrier.Set("foo", "bar1")
	carrier.Set("abc", "test")

	expected := PublishingCarrier{&amqp.Publishing{Headers: amqp.Table{
		"foo": "bar1",
		"abc": "test",
	}}}
	assert.Equal(t, carrier, expected)
}
