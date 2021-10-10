# eyewa-go-lib
Shared Go Lib for Eyewa's microservices.

# rabbitmq
This package provides an abstraction layer for the `github.com/streadway/amqp` pkg. It provides the capabilities required for publishing and consuming from RMQ within eyewa's microservices ecosystem. Pkg also supports having multiple instances of a consumer/publisher e.g:

- spinning up 1+ instances of the **catalogindexer** service to consume messages from the **catalogconsumer** queue
- spinning up 1+ instances of the **catalogconsumer** service to consume messages from the **eyewacatalog** queue
- etc


# How to use
The following variables should be injected in order to use this pkg

```go
// optional - used to identify what service is connected in RMQ's Admin UI.
"SERVICE_NAME"

// optional - if not set `amqp` will be used as protocol. if true `amqps`
"RABBITMQ_SECURED"

// required - RMQ credentials
"RABBITMQ_SERVER"
"RABBITMQ_AMQP_PORT"
"RABBITMQ_USERNAME"
"RABBITMQ_PASSWORD"

// the queues a service will be connecting to. At least one should be specified
// i.e either a service consumes or publishes or does both. if none is specified
// and pkg is included, it will yield an error.
"PUBLISHER_QUEUE_NAME" // queue service will be publishing to (optional)
"CONSUMER_QUEUE_NAME" // queue service will be consuming from (optional)

// optional - how many messages a consumer can consume at a go from RMQ.
// defaults to 5 if none is provided.
"QUEUE_PREFETCH_COUNT" 

// type of exchanges to use for a queue - fanout|direct|headers|topic
"RABBITMQ_PUBLISHER_EXCHANGE_TYPE" // required if PUBLISHER_QUEUE_NAME is provided
"RABBITMQ_CONSUMER_EXCHANGE_TYPE" // required if CONSUMER_QUEUE_NAME is provided
```

## Consuming from a Queue
Consuming from RMQ entails passing a callback func. For every message consumed from RMQ, the outcome is pushed to a callback func specified by the caller to act upon e.g persist event to datastore, or react to a failed message. On failed messages, such messages will be published to a deadletter queue for the queue. e.g `eyewacatalog` => `deadletter-eyewacatalog` etc.

There are two ways of consuming from RMQ using this client depending on the use case...

- using a Goroutine
- without a Goroutine

### Using a Goroutine
**Use Case**: I don't want consuming from a queue to be **blocking** to other operations. See code snippet below:

```go
	// consume from RMQ using a Goroutine
	// ideal if you dont want consuming to be blocking to other processes or Goroutines

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go broker.Client.Consume(config.Config.RabbitMQ.ConsumerQueueName, func(event *base.EyewaEvent, err error) {
		if err != nil {
			log.Error("Error consuming event",
				zap.String("queue", config.Config.RabbitMQ.ConsumerQueueName),
				zap.Error(err))
		} else {
			log.Info("Event received", zap.Any("event", event))
			// save to db, etc...
		}
	})

	fmt.Println("I will be called!")

	wg.Wait()
```

### Without a Goroutine
**Use Case**: I don't have any operations to perform besides consuming from a queue. See code snippet below

```go
	// consume from message broker
	// ideal if there is nothing else to do by caller other than consuming

	broker.Client.Consume(config.Config.RabbitMQ.ConsumerQueueName, func(event *base.EyewaEvent, err error) {
		if err != nil {
			log.Error("Error consuming event",
				zap.String("queue", config.Config.RabbitMQ.ConsumerQueueName),
				zap.Error(err))
		} else {
			log.Info("Event received", zap.Any("event", event))
			// save to db, etc...
		}
	})

	fmt.Println("I will NEVER be called!")
```

## Publishing to a Queue
Publishing to RMQ entails creating a publisher Goroutine with a callback func. This Goroutine lifecycle is not long lasting. Once it successfully publishes to RMQ, it dies off. So each publish call creates a seperate Goroutine to perform publishing.

```go 
	wg := new(sync.WaitGroup)
	wg.Add(1)

	event := &base.EyewaEvent{
		ID:        uuid.NewString(),
		Name:      "ProductCreated",
		EventType: "Products",
		Payload:   []byte(`{"test": "this is me testing"}`),
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	// publish to message broker
	// once the wg.Done() is called from within the Publish func this Goroutine seizes to exist.
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
```

## Publishing and Consuming
A service could require both publishing and consuming capabilities. In such cases, create 2 goroutines as seen below:

``` go
	// open connection to message broker
	broker, err := brokers.OpenConnection()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// ensure we wait for each Goroutine to fully perform its operations.
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

	...

	wg.Wait()
```