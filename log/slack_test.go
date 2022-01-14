package log

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestGetWebhook(t *testing.T) {
	t.Run("GetWebhook", func(t *testing.T) {
		webhook := GetWebhook()

		assert.NotNil(t, webhook)
	})

	t.Run("GetClient Singleton", func(t *testing.T) {
		webhook := GetWebhook()
		webhook2 := GetWebhook()

		assert.NotNil(t, webhook)
		assert.NotNil(t, webhook2)
		assert.Equal(t, webhook, webhook2)
	})
}

func TestSlackLogHook_Log(t *testing.T) {
	slackLogger = GetWebhook()

	err := slackLogger.Log("test message", zap.Field{})

	assert.Nil(t, err)
}
