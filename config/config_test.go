package config

import (
	"os"
	"testing"

	"github.com/eyewa/eyewa-go-lib/brokers"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	vars := map[string]string{
		"ENV":            "testing",
		"LOG_LEVEL":      "debug",
		"MESSAGE_BROKER": string(brokers.RabbitMQ),
	}
	for e, v := range vars {
		os.Setenv(e, v)
	}

	func(t *testing.T) {
		assert.Nil(t, initConfig())
	}(new(testing.T))

	exitVal := m.Run()

	for _, v := range envVars {
		os.Unsetenv(v)
	}

	os.Exit(exitVal)
}

func TestConfig(t *testing.T) {
	assert.NotZero(t, config)
	assert.Equal(t, "testing", config.Env)
	assert.Equal(t, "debug", config.LogLevel)
	assert.Equal(t, brokers.RabbitMQ, config.MessageBroker)
}

func TestGetConfigEnvVars(t *testing.T) {
	assert.NotEmpty(t, GetConfigEnvVars())
}
