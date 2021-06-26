package tracing

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/eyewa/eyewa-go-lib/tracing/exporter"
	"github.com/stretchr/testify/assert"
)

func TestLauncher(t *testing.T) {
	t.Run("connect to exporter on launch", func(tt *testing.T) {
		mockexp := new(exporter.Mock)
		ctx := context.Background()

		mockexp.On("Start", ctx).Return(nil)

		l := NewLauncher(mockexp)
		err := l.Launch()

		assert.Nil(t, err)
		assert.NoError(t, err)
		mockexp.AssertExpectations(t)
	})

	t.Run("connect to exporter on launch fail", func(tt *testing.T) {
		mockexp := new(exporter.Mock)
		ctx := context.Background()

		mockexp.On("Start", ctx).Return(errors.New("failed to connect"))

		l := NewLauncher(mockexp)
		err := l.Launch()

		assert.NotNil(t, err)
		assert.Error(t, err)
		mockexp.AssertExpectations(t)
	})

	t.Run("launcher should shutdown the exporter on shutdown", func(tt *testing.T) {
		mockexp := new(exporter.Mock)
		ctx := context.Background()

		mockexp.On("Start", ctx).Return(nil)
		mockexp.On("Shutdown", ctx).Return(nil)

		l := NewLauncher(mockexp)
		err := l.Launch()
		if err != nil {
			t.Fatalf("expected launcher to launch: %v", err)
		}

		err = l.Shutdown()

		mockexp.AssertCalled(t, "Shutdown", l.ctx)
		mockexp.AssertNumberOfCalls(t, "Shutdown", 1)
		assert.Nil(t, err)
		assert.NoError(t, err)
		mockexp.AssertExpectations(t)
	})

	t.Run("launcher should return an error when failing to start exporter", func(tt *testing.T) {
		mockexp := new(exporter.Mock)
		ctx := context.Background()

		mockexp.On("Start", ctx).Return(errors.New("failed to connect."))

		l := NewLauncher(mockexp)
		err := l.Launch()
		fmt.Println(err)
		assert.NotNil(t, err)
		assert.Error(t, err)
		mockexp.AssertExpectations(t)
	})

	t.Run("launcher should return an error when failing to shutdown exporter", func(tt *testing.T) {
		mockexp := new(exporter.Mock)
		ctx := context.Background()

		mockexp.On("Start", ctx).Return(nil)
		mockexp.On("Shutdown", ctx).Return(errors.New("connection timed out."))

		l := NewLauncher(mockexp)
		err := l.Launch()
		if err != nil {
			t.Fatalf("expected launcher to launch: %v", err)
		}

		err = l.Shutdown()

		assert.NotNil(t, err)
		assert.Error(t, err)
		mockexp.AssertExpectations(t)
	})
}

func TestLaunch(t *testing.T) {
	t.Run("should launch with no error", func(tt *testing.T) {
		mockexp := new(exporter.Mock)

		mockexp.AssertNotCalled(t, "Start")
		_, err := Launch()

		assert.Nil(t, err)
		assert.NoError(t, err)
		mockexp.AssertExpectations(t)
	})

	t.Run("should return error when fails to launch", func(tt *testing.T) {
		mockexp := new(exporter.Mock)

		mockexp.AssertNotCalled(t, "Start")
		_, err := Launch()

		assert.Nil(t, err)
		assert.NoError(t, err)
		mockexp.AssertExpectations(t)
	})
}

func TestConfiguration(t *testing.T) {

}
