package tracing

import (
	"os"
	"testing"

	"github.com/eyewa/eyewa-go-lib/log"
	"github.com/stretchr/testify/assert"
)

func setup() func() {
	os.Clearenv()
	os.Setenv("LOG_LEVEL", "info")
	log.SetLogLevel()
	config = Config{}
	return func() {
		os.Clearenv()
	}
}

func TestLaunchWithoutServiceName(t *testing.T) {
	teardown := setup()
	defer teardown()

	os.Setenv("TRACING_EXPORTER_ENDPOINT", "fake-endpoint.test")

	shutdown, err := Launch()
	defer func() {
		shutdownErr := shutdown()
		assert.Nil(t, shutdownErr)
		assert.NoError(t, shutdownErr)
	}()

	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.NotZero(t, err)
}

func TestLaunchWithoutEndpoint(t *testing.T) {
	teardown := setup()
	defer teardown()

	os.Setenv("SERVICE_NAME", "test-service")

	shutdown, err := Launch()

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
	teardown := setup()
	defer teardown()

	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("TRACING_EXPORTER_ENDPOINT", "fake-endpoint.test")

	shutdown, err := Launch()
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
	teardown := setup()
	defer teardown()

	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("TRACING_EXPORTER_ENDPOINT", "fake-endpoint.test")

	shutdown, err := Launch()
	defer func() {
		err = shutdown()
		assert.Nil(t, err)
		assert.NoError(t, err)
	}()

	assert.Nil(t, err)
	assert.NoError(t, err)
	assert.Zero(t, err)

	// relaunch
	shutdown, err = Launch()
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
	teardown := setup()
	defer teardown()

	os.Setenv("SERVICE_NAME", "testing")
	os.Setenv("TRACING_SECURE_EXPORTER", "fake-endpoint.test")

	shutdown, err := Launch()
	defer func() {
		err = shutdown()
		assert.Nil(t, err)
		assert.NoError(t, err)
	}()

	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.NotZero(t, err)
}

// func TestBlockingExporter(t *testing.T) {
// 	teardown := setup()
// 	defer teardown()

// 	os.Setenv("SERVICE_NAME", "testing")
// 	os.Setenv("TRACING_EXPORTER_ENDPOINT", "fake-endpoint.test")
// 	os.Setenv("TRACING_BLOCK_EXPORTER", "true")

// 	shutdown, err := Launch()
// 	defer func() {
// 		err = shutdown()
// 		assert.Nil(t, err)
// 		assert.NoError(t, err)
// 	}()

// 	assert.Nil(t, err)
// 	assert.NoError(t, err)
// 	assert.Zero(t, err)
// }

func TestSecureExporter(t *testing.T) {
	teardown := setup()
	defer teardown()

	os.Setenv("SERVICE_NAME", "testing")
	os.Setenv("TRACING_EXPORTER_ENDPOINT", "fake-endpoint.test")
	os.Setenv("TRACING_SECURE_EXPORTER", "true")

	shutdown, err := Launch()
	defer func() {
		err = shutdown()
		assert.Nil(t, err)
		assert.NoError(t, err)
	}()

	assert.Nil(t, err)
	assert.NoError(t, err)
	assert.Zero(t, err)
}
