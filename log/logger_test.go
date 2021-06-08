package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSetLogLevel(t *testing.T) {
	os.Setenv("LOG_LEVEL", "info")
	SetLogLevel()

	check := logger.Check(zap.DebugLevel, "sss")
	assert.Nil(t, check)

	os.Setenv("LOG_LEVEL", "debug")
	SetLogLevel()
	check = logger.Check(zap.DebugLevel, "sss")
	assert.Equal(t, zap.DebugLevel, check.Level)
}
