package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrimitivePointerConversions(t *testing.T) {
	assert.IsType(t, new(int), ConvertIntToPointer(2))
	assert.IsType(t, new(string), ConvertStringToPointer("blah"))
	assert.IsType(t, new(float64), ConvertFloat64ToPointer(455))
}
