# eyewa-go-lib

Shared Go Lib for Eyewa's microservices.

# amqp

This pkg provides `tracing` decorators for the `github.com/streadway/amqp` pkg. It decorates a `amqp.Publishing` and `amqp.Delivery` message by starting a trace whereby it's headers are extracted/injected to enable the propagation of trace context across services.

This pkg is mostly used internally within the `rabbitmq` implementation of the `brokers` pkg.


# How to use

### Starting A Delivery Trace

```go
package main

import (
    rmqtrace "github.com/eyewa/eyewa-go-lib/tracing/amqp"
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
    rmqtrace "github.com/eyewa/eyewa-go-lib/tracing/amqp"
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
