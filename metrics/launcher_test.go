package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewLauncher(t *testing.T) {
	option := ExportOption{
		CollectPeriod: 10 * time.Second,
	}

	ml, err := NewLauncher(option)
	assert.Nil(t, err)
	assert.NotNil(t, ml)
}
