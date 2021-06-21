# eyewa-go-lib
Shared Go Lib for Eyewa's microservices.

# brokers
This package provides an abstraction layer for any underlying third party pkgs we utilize in publishing or consuming messages from a message broker of choice. Currently there are **provisional clients** for SQS, RabbitMQ and Kafka. Concrete implementations will be added from time to time to support requirements and any client capabilities we require.

A client can either be a **Consumer**, a **Publisher** or **both** - there are contracts in place to cater for such scenarios. On calling the `OpenConnection`, a client will be regarded as requiring both capabilities. If otherwise, there are direct client calls for initiating either a consumer/publisher for any client of choice.

# How to use

```go
package myservice

import (
	"sync"
	"time"

	"github.com/eyewa/catalogconsumer-service/config"
	"github.com/eyewa/eyewa-go-lib/base"
	"github.com/eyewa/eyewa-go-lib/brokers"
	"github.com/eyewa/eyewa-go-lib/log"
	"github.com/eyewa/eyewa-go-lib/uuid"
	"go.uber.org/zap"
)

func main() {
	// this should be injected and not hardcoded.
	os.Setenv("MESSAGE_BROKER", "rabbitmq") 

	err := config.Init()
	if err != nil {
		log.Error(err.Error())
	}

	// open connection to message broker
	broker, err := brokers.OpenConnection()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// ensure we wait for each go routine to fully perform its operations.
	wg := new(sync.WaitGroup)
	wg.Add(2)

	event := &base.EyewaEvent{
		ID:        uuid.NewString(),
		Name:      "ProductCreated",
		EventType: "Products",
		Payload:   []byte(`{"test": "this is me testing"}`),
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	// publish to message broker
	go broker.Client.Publish(config.Config.RabbitMQ.PublisherQueueName, event, func(event *base.EyewaEvent, err error) {
		log.Debug("Publishing event", zap.String("event", event.ID))

		if err != nil {
			log.Error("Failed publishing event",
				zap.String("queue", config.Config.RabbitMQ.PublisherQueueName),
				zap.Error(err))
		} else {
			log.Info("Event was published successfully", zap.String("event", event.ID))
		}
	}, wg)

	// consume from message broker
	go broker.Client.Consume(config.Config.RabbitMQ.ConsumerQueueName, func(event *base.EyewaEvent, err error) {
		log.Debug("Error consuming events",
			zap.String("queue", config.Config.RabbitMQ.ConsumerQueueName),
			zap.Error(err))

		log.Debug("Event received", zap.Any("event", event))
	})

	wg.Wait()


	// inititate consumer client only
	consumer := brokers.NewConsumerClient(brokers.RabbitMQ)
	err := consumer.Client.Connect()
	if err != nil {
		log.Error("Cannot connect to broker.", zap.Error(err), zap.Any("client", consumer.Client))
	}

	// initiate publisher client only
	publisher := brokers.NewPublisherClient(brokers.RabbitMQ)
	err := publisher.Client.Connect()
	if err != nil {
		log.Error("Cannot connect to broker.", zap.Error(err), zap.Any("client", publisher.Client))
	}
}
```
