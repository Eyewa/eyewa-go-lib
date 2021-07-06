package tracing

import (
	"context"
	"os"
	"testing"

	"github.com/eyewa/eyewa-go-lib/log"
	"github.com/stretchr/testify/assert"
)

func TestNewResource(t *testing.T) {
	os.Clearenv()
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("LOG_LEVEL", "debug")
	log.SetLogLevel()
	_, err := initConfig()
	assert.Nil(t, err)

	res, _ := newResource(context.Background())

	assert.NotNil(t, res)
	assert.NotNil(t, res.Attributes())
	assert.Greater(t, res.Len(), 0)
}
