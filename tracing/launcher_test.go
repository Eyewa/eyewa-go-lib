package tracing_test

import (
	"os"
	"testing"

	"github.com/eyewa/eyewa-go-lib/log"
	"github.com/eyewa/eyewa-go-lib/tracing"
	"github.com/stretchr/testify/assert"
)

func reset() {
	os.Clearenv()
	os.Setenv("LOG_LEVEL", "debug")
	log.SetLogLevel()
}

func TestLaunchWithoutServiceName(t *testing.T) {
	reset()
	os.Setenv("TRACING_EXPORTER_ENDPOINT", "fake-endpoint.test")

	shutdown, err := tracing.Launch()
	defer func() {
		err = shutdown()
		assert.Nil(t, err)
		assert.NoError(t, err)
	}()

	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.NotZero(t, err)

}

func TestLaunchWithoutEndpoint(t *testing.T) {
	reset()
	os.Setenv("SERVICE_NAME", "test-service")

	shutdown, err := tracing.Launch()
	defer func() {
		err = shutdown()
		assert.Nil(t, err)
		assert.NoError(t, err)
	}()

	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.NotZero(t, err)
}

func TestLaunchShutdown(t *testing.T) {
	reset()
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("TRACING_EXPORTER_ENDPOINT", "fake-endpoint.test")

	shutdown, err := tracing.Launch()
	defer func() {
		err = shutdown()
		assert.Nil(t, err)
		assert.NoError(t, err)
	}()

	assert.Nil(t, err)
	assert.NoError(t, err)
	assert.Zero(t, err)
}

func TestShutdownRelaunch(t *testing.T) {
	reset()
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("TRACING_EXPORTER_ENDPOINT", "fake-endpoint.test")

	shutdown, err := tracing.Launch()
	defer func() {
		err = shutdown()
		assert.Nil(t, err)
		assert.NoError(t, err)
	}()

	assert.Nil(t, err)
	assert.NoError(t, err)
	assert.Zero(t, err)

	// relaunch
	shutdown, err = tracing.Launch()
	defer func() {
		err = shutdown()
		assert.Nil(t, err)
		assert.NoError(t, err)
	}()

	assert.Nil(t, err)
	assert.NoError(t, err)
	assert.Zero(t, err)
}

func TestInvalidEnvConfig(t *testing.T) {
	reset()
	os.Setenv("SERVICE_NAME", "testing")
	os.Setenv("TRACING_SECURE_EXPORTER", "fake-endpoint.test")

	shutdown, err := tracing.Launch()
	defer func() {
		err = shutdown()
		assert.Nil(t, err)
		assert.NoError(t, err)
	}()

	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.NotZero(t, err)
}

func TestBlockingExporter(t *testing.T) {
	reset()
	os.Setenv("SERVICE_NAME", "testing")
	os.Setenv("TRACING_EXPORTER_ENDPOINT", "fake-endpoint.test")
	os.Setenv("TRACING_BLOCK_EXPORTER", "true")

	shutdown, err := tracing.Launch()
	defer func() {
		err = shutdown()
		assert.Nil(t, err)
		assert.NoError(t, err)
	}()

	assert.Nil(t, err)
	assert.NoError(t, err)
	assert.Zero(t, err)
}

func TestSecureExporter(t *testing.T) {
	reset()
	os.Setenv("SERVICE_NAME", "testing")
	os.Setenv("TRACING_EXPORTER_ENDPOINT", "fake-endpoint.test")
	os.Setenv("TRACING_SECURE_EXPORTER", "true")

	shutdown, err := tracing.Launch()
	defer func() {
		err = shutdown()
		assert.Nil(t, err)
		assert.NoError(t, err)
	}()

	assert.Nil(t, err)
	assert.NoError(t, err)
	assert.Zero(t, err)
}
