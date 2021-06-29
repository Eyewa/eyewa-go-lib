# eyewa-go-lib

Shared Go Lib for Eyewa's microservices.

## tracing

This package configures Open Telemetry as the global tracing provider. It configures an endpoint to where all traces end up using the `TRACING_EXPORTER_ENDPOINT` env variable. A user of this package is able to view traces on Grafana Tempo once a trace has been exported.

</br>

## How To Use

- Set the `SERVICE_NAME` environmental variable.
- Set the `TRACING_EXPORTER_ENDPOINT` environmental variable.
- launch to connect to the open telemetry collector.
- Add a GRPC interceptor to the GRPC server/client.

</br>

### Environmental Variables

```go
SERVICE_NAME // Name of the service/application. #Required
TRACING_EXPORTER_ENDPOINT // The endpoint that spans get exported to. #Required
TRACING_BLOCK_EXPORTER // Exporter initiates a blocking request to an endpoint | #Optional | bool
TRACING_SECURE_EXPORTER // Exporter connects with TLS secure connection. | #Optional | bool
```

</br>

### GRPC Server Tracing

```go
package exampleservice

import (
 "net"
 "os"

 "github.com/eyewa/eyewa-go-lib/log"
 "github.com/eyewa/eyewa-go-lib/tracing"
 "github.com/eyewa/exampleservice/api"
 "github.com/eyewa/exampleservice/config"
 "google.golang.org/grpc"
)

func main() {
 // this should be injected and not hardcoded.
 os.Setenv("SERVICE_NAME", "exampleservice")
 os.Setenv("TRACING_EXPORTER_ENDPOINT", "open-telemetry.collector.endpoint")
 os.Setenv("GRPC_SERVER_PORT", "7777")

 err := config.Init()
 if err != nil {
  log.Error(err.Error())
  return
 }

 // launch tracing to open a connection to
 // a tracing backend.
 shutdown, err := tracing.Launch()
 defer shutdown()
 if err != nil {
  log.Error(err.Error())
  return
 }


 // setup the service grpc server as normal.
 port := os.Getenv("GRPC_SERVER_PORT")
 lis, err := net.Listen("tcp", port)
 defer lis.Close()
 if err != nil {
  log.Error(err.Error())
  return
 }

 // inject tracing interceptors.
 s := grpc.NewServer(
  tracing.UnaryServerTraceInterceptor(),
  tracing.StreamServerTraceInterceptor(),
 )

 // register the server and start serving grpc requests.
 api.RegisterHelloServiceServer(s, &server{})
 if err := s.Serve(lis); err != nil {
  log.Error(err.Error())
  return
 }

}

```

</br>

### RabbitMQ Consumer Tracing

```go
// This example declares a durable Exchange, an ephemeral (auto-delete) Queue,
// binds the Queue to the Exchange with a binding key, and consumes every
// message published to that Exchange with that routing key.
//
package main

import (
    "flag"
    "fmt"
    "github.com/streadway/amqp"
    "log"
    "time"
    "github.com/eyewa/eyewa-go-lib/tracing"
)

var (
 uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
 exchange     = flag.String("exchange", "test-exchange", "Durable, non-auto-deleted AMQP exchange name")
 exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
 queue        = flag.String("queue", "test-queue", "Ephemeral AMQP queue name")
 bindingKey   = flag.String("key", "test-key", "AMQP binding key")
 consumerTag  = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
 lifetime     = flag.Duration("lifetime", 5*time.Second, "lifetime of process before shutdown (0s=infinite)")
)

func init() {
 flag.Parse()
}

func main() {
    c, err := NewConsumer(*uri, *exchange, *exchangeType, *queue, *bindingKey, *consumerTag)
    if err != nil {
    log.Fatalf("%s", err)
    }

    if *lifetime > 0 {
    time.Sleep(*lifetime)
    } else {
    select {}
    }


    if err := c.Shutdown(); err != nil {
    log.Error("error during shutdown: %s", err)
    return
    }
}

type Consumer struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    tag     string
    done    chan error
}

func NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Consumer, error) {
    c := &Consumer{
    conn:    nil,
    channel: nil,
    tag:     ctag,
    done:    make(chan error),
    }

    var err error

    c.conn, err = amqp.Dial(amqpURI)
    if err != nil {
    return nil, fmt.Errorf("Dial: %s", err)
    }

    go func() {
    fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
    }()

    log.Printf("got Connection, getting Channel")
    c.channel, err = c.conn.Channel()
    if err != nil {
    return nil, fmt.Errorf("Channel: %s", err)
    }

    if err = c.channel.ExchangeDeclare(
    exchange,     // name of the exchange
    exchangeType, // type
    true,         // durable
    false,        // delete when complete
    false,        // internal
    false,        // noWait
    nil,          // arguments
    ); err != nil {
    return nil, fmt.Errorf("Exchange Declare: %s", err)
    }

    queue, err := c.channel.QueueDeclare(
    queueName, // name of the queue
    true,      // durable
    false,     // delete when unused
    false,     // exclusive
    false,     // noWait
    nil,       // arguments
    )
    if err != nil {
    return nil, fmt.Errorf("Queue Declare: %s", err)
    }

    queue.Name, queue.Messages, queue.Consumers, key)

    if err = c.channel.QueueBind(
    queue.Name, // name of the queue
    key,        // bindingKey
    exchange,   // sourceExchange
    false,      // noWait
    nil,        // arguments
    ); err != nil {
    return nil, fmt.Errorf("Queue Bind: %s", err)
    }

    deliveries, err := c.channel.Consume(
    queue.Name, // name
    c.tag,      // consumerTag,
    false,      // noAck
    false,      // exclusive
    false,      // noLocal
    false,      // noWait
    nil,        // arguments
    )
    if err != nil {
    return nil, fmt.Errorf("Queue Consume: %s", err)
    }

    go handle(deliveries, c.done)

    return c, nil
}

func (c *Consumer) Shutdown() error {
    if err := c.channel.Cancel(c.tag, true); err != nil {
    return fmt.Errorf("Consumer cancel failed: %s", err)
    }

    if err := c.conn.Close(); err != nil {
    return fmt.Errorf("AMQP connection close error: %s", err)
    }

    return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
    for d := range deliveries {
    // Extract tracing info from message
    ctx := otel.GetTextMapPropagator().Extract(
            context.Background(), 
            tracing.NewRabbitMQDeliveryCarrier(msg)
        )

    d.Ack(false)
    }
    log.Printf("handle: deliveries channel closed")
    done <- nil
}

```

</br>
