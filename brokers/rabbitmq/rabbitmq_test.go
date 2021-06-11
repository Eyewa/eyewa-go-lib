package rabbitmq

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionConfig(t *testing.T) {
	vars := map[string]string{
		"RABBITMQ_SERVER":    "localhost",
		"RABBITMQ_AMQP_PORT": "11111",
		"RABBITMQ_USERNAME":  "bleh",
		"RABBITMQ_PASSWORD":  "blah",
	}
	for e, v := range vars {
		os.Setenv(e, v)
	}

	cfg, str, err := initConfig()
	assert.NotZero(t, cfg)
	assert.Nil(t, err)
	assert.True(t, true)
	assert.Equal(t, str, fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Username, config.Password, config.Server, config.AmqpPort))

	os.Clearenv()
}

func TestConnectionConfigNotSet(t *testing.T) {
	config = Config{}
	cfg, _, err := initConfig()
	assert.Zero(t, cfg)
	assert.Nil(t, err)
	assert.True(t, true)
}
