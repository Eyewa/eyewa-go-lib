package metrics

import (
	"github.com/ory/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLauncher(t *testing.T) {
	viper.Set("METRIC_COLLECT_PERIOD", "5s")

	ml, err := NewLauncher()
	assert.Nil(t, err)
	assert.NotNil(t, ml)
}
