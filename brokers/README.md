# eyewa-go-lib
Shared Go Lib for Eyewa's microservices.

# brokers
This package provides an abstraction layer for any underlying third party pkgs we utilize in publishing or consuming messages from a message broker of choice. Currently there are **provisional clients** for SQS, RabbitMQ and Kafka. Concrete implementations will be added from time to time to support requirements and any client capabilities we require.

A client can either be a **Consumer**, a **Publisher** or **both** - there are contracts in place to cater for such scenarios. On calling the `OpenConnection`, a client will be regarded as requiring both capabilities. If otherwise, there are direct client calls for initiating either a consumer/publisher for any client of choice.


```go
package demo

"github.com/eyewa/eyewa-go-lib/broker"

func main() {
  os.Setenv("MESSAGE_BROKER", "rabbitmq") // this should be injected and not hardcoded.

  // BROKER CLIENT
  broker, err := brokers.OpenConnection()
  if err != nil {
    log.Error(err.Error())
  }
  _ = broker.Client.Consume("events")
  _ = broker.Client.Publish("events")
  _ = broker.Client.CloseConnection()
  log.Debug("I am using", zap.String("broker", string(broker.Type)))


  // CONSUMER CLIENT
  consumer := brokers.NewConsumerClient(brokers.RabbitMQ)
  err := consumer.Client.Connect()
  if err != nil {
    log.Error("Cannot connect to broker.", zap.Error(err), zap.Any("client", consumer.Client))
  }

  // PUBLISHER CLIENT
  publisher := brokers.NewPublisherClient(brokers.RabbitMQ)
  err := publisher.Client.Connect()
  if err != nil {
    log.Error("Cannot connect to broker.", zap.Error(err), zap.Any("client", publisher.Client))
  }
}
```
