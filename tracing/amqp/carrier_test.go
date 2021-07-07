package amqp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublishingHeaderCarrierGet(t *testing.T) {
	testCases := []struct {
		name     string
		carrier  HeaderCarrier
		key      string
		expected string
	}{
		{
			name: "header exists",
			carrier: HeaderCarrier{
				"foo": "bar",
			},
			key:      "foo",
			expected: "bar",
		},
		{
			name:     "header does not exists",
			carrier:  HeaderCarrier{},
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
	carrier := HeaderCarrier{
		"foo": "bar",
	}

	carrier.Set("foo", "bar1")
	carrier.Set("abc", "test")
	carrier.Set("hello", "world")

	expected := HeaderCarrier{
		"foo":   "bar1",
		"abc":   "test",
		"hello": "world",
	}
	assert.Equal(t, carrier, expected)
}

func TestPublishingHeaderCarrierKeys(t *testing.T) {
	testCases := []struct {
		name     string
		carrier  HeaderCarrier
		expected []string
	}{
		{
			name: "one header",
			carrier: HeaderCarrier{
				"foo": "bar1",
			},
			expected: []string{"foo"},
		},
		{
			name:     "no headers",
			carrier:  HeaderCarrier{},
			expected: []string{},
		},
		{
			name: "multiple headers",
			carrier: HeaderCarrier{
				"foo":   "bar1",
				"abc":   "test",
				"hello": "world",
			},
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
