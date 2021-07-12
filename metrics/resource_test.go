package metrics

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResource(t *testing.T) {
	os.Clearenv()
	os.Setenv("SERVICE_NAME", "test-service")

	_, err := initConfig()
	assert.Nil(t, err)

	option, err := initConfig()
	assert.Nil(t, err)

	res, _ := newResource(option)

	assert.NotNil(t, res)
	assert.NotNil(t, res.Attributes())
	assert.Greater(t, res.Len(), 0)
}
