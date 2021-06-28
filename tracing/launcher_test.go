package tracing_test

import (
	"os"
	"testing"

	"github.com/eyewa/eyewa-go-lib/errors"
	"github.com/eyewa/eyewa-go-lib/tracing"
	"github.com/stretchr/testify/assert"
)

func TestLaunchWithServiceName(t *testing.T) {
	os.Clearenv()
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("EXPORTER_ENDPOINT", "test-endpoint")

	shutdown, err := tracing.Launch()

	assert.Nil(t, err)
	assert.NoError(t, err)

	err = shutdown()
	assert.Nil(t, err)
	assert.NoError(t, err)
}

func TestLaunchWithoutServiceName(t *testing.T) {
	os.Clearenv()
	os.Setenv("EXPORTER_ENDPOINT", "test-endpoint")
	shutdown, err := tracing.Launch()
	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errors.ErrorNoServiceNameSpecified)

	err = shutdown()
	assert.Nil(t, err)
	assert.NoError(t, err)
}

func TestLaunchShutdown(t *testing.T) {
	os.Clearenv()
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("EXPORTER_ENDPOINT", "test-endpoint")
	shutdown, err := tracing.Launch()
	assert.Nil(t, err)
	assert.NoError(t, err)

	err = shutdown()
	assert.Nil(t, err)
	assert.NoError(t, err)
}

func TestLaunchWithoutEndpoint(t *testing.T) {
	os.Clearenv()
	os.Setenv("SERVICE_NAME", "test-service")

	shutdown, err := tracing.Launch()
	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errors.ErrorNoExporterEndpointSpecified)

	err = shutdown()
	assert.Nil(t, err)
	assert.NoError(t, err)
}
func TestShutdownRelaunch(t *testing.T) {
	os.Clearenv()
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("EXPORTER_ENDPOINT", "123")

	shutdown, err := tracing.Launch()
	assert.Nil(t, err)
	assert.NoError(t, err)

	err = shutdown()
	assert.Nil(t, err)
	assert.NoError(t, err)

	// relaunch
	shutdown, err = tracing.Launch()
	assert.Nil(t, err)
	assert.NoError(t, err)

	err = shutdown()
	assert.Nil(t, err)
	assert.NoError(t, err)
}

func TestInvalidEnvConfig(t *testing.T) {
	os.Setenv("SERVICE_NAME", "testing")
	os.Setenv("EXPORTER_SECURE", "none")
	_, err := tracing.Launch()

	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.NotZero(t, err)
}

func TestBlockingExporterFail(t *testing.T) {
	os.Setenv("SERVICE_NAME", "testing")
	os.Setenv("EXPORTER_BLOCKING", "true")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("EXPORTER_ENDPOINT", "fake-endpoint.test")
	_, err := tracing.Launch()

	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.NotZero(t, err)
}
