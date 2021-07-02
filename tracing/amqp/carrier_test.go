package amqp

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestNewPublishingHeaderCarrier(t *testing.T) {
	headers := amqp.Table{"foo": "bar"}
	pub := &amqp.Publishing{Headers: headers}
	c := NewPublishingHeaderCarrier(pub)

	assert.NotNil(t, c)
	assert.NotZero(t, c)
	assert.Equal(t, len(headers), len(c.Keys()))
}

func TestPublishingHeaderCarrierGet(t *testing.T) {
	testCases := []struct {
		name     string
		carrier  PublishingHeaderCarrier
		key      string
		expected string
	}{
		{
			name: "header exists",
			carrier: PublishingHeaderCarrier{publishing: &amqp.Publishing{Headers: amqp.Table{
				"foo": "bar",
			}}},
			key:      "foo",
			expected: "bar",
		},
		{
			name:     "header does not exists",
			carrier:  PublishingHeaderCarrier{publishing: &amqp.Publishing{Headers: amqp.Table{}}},
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

func TestPublishingHeaderCarrierSet(t *testing.T) {
	pub := amqp.Publishing{Headers: amqp.Table{
		"foo": "bar",
	}}
	carrier := PublishingHeaderCarrier{publishing: &pub}

	carrier.Set("foo", "bar1")
	carrier.Set("abc", "test")
	carrier.Set("hello", "world")

	expected := PublishingHeaderCarrier{&amqp.Publishing{Headers: amqp.Table{
		"foo":   "bar1",
		"abc":   "test",
		"hello": "world",
	}}}
	assert.Equal(t, carrier, expected)
}

func TestPublishingHeaderCarrierKeys(t *testing.T) {
	testCases := []struct {
		name     string
		carrier  PublishingHeaderCarrier
		expected []string
	}{
		{
			name: "one header",
			carrier: PublishingHeaderCarrier{&amqp.Publishing{Headers: amqp.Table{
				"foo": "bar1",
			}}},
			expected: []string{"foo"},
		},
		{
			name: "no headers",
			carrier: PublishingHeaderCarrier{publishing: &amqp.Publishing{
				Headers: amqp.Table{}}},
			expected: []string{},
		},
		{
			name: "multiple headers",
			carrier: PublishingHeaderCarrier{&amqp.Publishing{Headers: amqp.Table{
				"foo":   "bar1",
				"abc":   "test",
				"hello": "world",
			}}},
			expected: []string{"foo", "abc", "hello"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.carrier.Keys()
			assert.ElementsMatch(t, tc.expected, result)
		})
	}
}

func TestNewDeliveryCarrier(t *testing.T) {
	headers := amqp.Table{"foo": "bar"}
	pub := amqp.Delivery{Headers: headers}
	c := NewDeliveryCarrier(pub)

	assert.NotNil(t, c)
	assert.NotZero(t, c)
	assert.Equal(t, len(headers), len(c.Keys()))
}

func TestDeliveryCarrierGet(t *testing.T) {
	testCases := []struct {
		name     string
		carrier  DeliveryCarrier
		key      string
		expected string
	}{
		{
			name: "header exists",
			carrier: DeliveryCarrier{delivery: amqp.Delivery{Headers: amqp.Table{
				"foo": "bar",
			}}},
			key:      "foo",
			expected: "bar",
		},
		{
			name:     "header does not exists",
			carrier:  DeliveryCarrier{delivery: amqp.Delivery{Headers: amqp.Table{}}},
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

func TestDeliveryHeaderCarrierSet(t *testing.T) {
	d := amqp.Delivery{Headers: amqp.Table{
		"foo": "bar",
	}}
	carrier := DeliveryCarrier{delivery: d}

	carrier.Set("foo", "bar1")
	carrier.Set("abc", "test")
	carrier.Set("hello", "world")

	expected := DeliveryCarrier{amqp.Delivery{Headers: amqp.Table{
		"foo":   "bar1",
		"abc":   "test",
		"hello": "world",
	}}}
	assert.Equal(t, carrier, expected)
}

func TestDeliveryHeaderCarrierKeys(t *testing.T) {
	testCases := []struct {
		name     string
		carrier  DeliveryCarrier
		expected []string
	}{
		{
			name: "one header",
			carrier: DeliveryCarrier{amqp.Delivery{Headers: amqp.Table{
				"foo": "bar1",
			}}},
			expected: []string{"foo"},
		},
		{
			name: "no headers",
			carrier: DeliveryCarrier{amqp.Delivery{
				Headers: amqp.Table{}}},
			expected: []string{},
		},
		{
			name: "multiple headers",
			carrier: DeliveryCarrier{amqp.Delivery{Headers: amqp.Table{
				"foo":   "bar1",
				"abc":   "test",
				"hello": "world",
			}}},
			expected: []string{"foo", "abc", "hello"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.carrier.Keys()
			assert.ElementsMatch(t, tc.expected, result)
		})
	}
}
