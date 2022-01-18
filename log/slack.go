package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strconv"
	"sync"

	"github.com/slack-go/slack"
)

type SlackLogHook interface {
	Log(message string, fields ...zap.Field) error
}

type slackLogHook struct {
	*slack.WebhookMessage
}

var (
	webhookOnce sync.Once
	webhook     *slackLogHook
)

func GetWebhook() *slackLogHook {
	webhookOnce.Do(func() {
		webhook = NewWebhook()
	})

	return webhook
}

func NewWebhook() *slackLogHook {
	return &slackLogHook{&slack.WebhookMessage{}}
}

func (sh *slackLogHook) Log(message string, fields ...zap.Field) error {
	if !(os.Getenv("LOG_ENABLE_SLACK") == "true") {
		return nil
	}

	payload := sh.WebhookMessage

	attachment := slack.Attachment{
		Pretext:    message,
		FooterIcon: "https://cdn.eyewa.com/media/favicon/default/Favicon_32x32.png",
		Color:      "danger",
		Footer:     "Errors",
	}

	var attachmentFields []slack.AttachmentField

	attachmentFields = append(attachmentFields, slack.AttachmentField{
		Title: "Environment",
		Value: os.Getenv("ENV"),
		Short: true,
	})

	for _, field := range fields {
		value := field.String

		switch field.Type {
		case zapcore.Int16Type:
		case zapcore.Int64Type:
			value = strconv.FormatInt(field.Integer, 10)
		case zapcore.ErrorType:
			err := field.Interface.(error)
			value = err.Error()
		}

		attachmentFields = append(attachmentFields, slack.AttachmentField{
			Title: field.Key,
			Value: value,
			Short: len(value) < 20,
		})
	}

	attachment.Fields = attachmentFields
	payload.Attachments = []slack.Attachment{attachment}

	return slack.PostWebhook(os.Getenv("SLACK_WEBHOOK_URL"), payload)
}
