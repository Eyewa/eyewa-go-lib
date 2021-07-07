package metrics

import (
	"testing"

	"github.com/ory/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewLauncher(t *testing.T) {
	viper.Set("METRIC_COLLECT_PERIOD", "5s")

	ml, err := newLauncher()
	assert.Nil(t, err)
	assert.NotNil(t, ml)
}
