# eyewa-go-lib

Shared Go Lib for Eyewa's microservices.

## amqp

This pkg provides `tracing` support for the `github.com/streadway/amqp` pkg. This pkg is mostly used internally within the `rabbitmq` implementation of the `brokers` pkg.

### Publishing trace process

1. Check if there is an existing trace context in the `amqp.Publishing` to use as a parent trace context by extracting its headers from the carrier.

2. Start a new span with attributes relating to `amqp` and an `amqp.Publishing`.

3. Inject the new context received from the new span into the `amqp.Publishing`

### Publishing trace process

1. Check if there is an existing trace context in the `amqp.Publishing` to use as a parent trace context by extracting its headers from the carrier.

2. Start a new span with attributes relating to `amqp` and an `amqp.Publishing`.

3. Inject the new context received from the new span into the `amqp.Publishing`



## How to use

### Starting A Delivery Trace

```go
package main

import (
    amqptracing "github.com/eyewa/eyewa-go-lib/tracing/amqp"
)

var (
    queue = "mycool.queue"
)

func main(){
    //... initialise exchange, queue etc...

    // attempt to consume events from broker.
    msgs, err := channel.Consume(queue, getNameForChannel(queue), false, false, false, false, nil)

    for d := range msgs {

        // start a span and extract the context
        ctx, endSpan = amqptracing.StartDeliverySpan(ctx, d)
        defer endSpan()

        // processMessage using the context
        processMessage(ctx, d)

    }
}
```

### Starting A Publishing Trace

```go
package main

import (
    amqptracing "github.com/eyewa/eyewa-go-lib/tracing/amqp"
)

var (
    queue = "mycool.queue"
)

func main(){
    //... initialise exchange, queue etc...

   msg := amqp.Publishing{
        ContentType:  "application/json",
        Body:         eventJSON,
        DeliveryMode: amqp.Persistent,
    }

    // start tracing the publishing span
    ctx, endSpan := amqptracing.StartPublishingSpan(ctx, msg)
    defer endSpan()

    // attempt to publish event
    err = channel.Publish("", queue, false, false, msg)
}
```
