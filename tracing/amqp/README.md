# eyewa-go-lib

Shared Go Lib for Eyewa's microservices.

# amqp

This pkg provides `tracing` decorators for the `github.com/streadway/amqp` pkg. It decorates the `Consume` and `Publish` methods on a `amqp.Channel`. Every `amqp.Publishing` and `amqp.Delivery` message has it's headers modified to enable the propagation of a trace.

This pkg is mostly used internally within `rabbitmq` implementation in the `brokers` pkg.


# How to use

### Trace All Deliveries

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
    consumeChannel := rmqtrace.WrapConsume(channel.Consume)
    msgs, err := consumeChannel(queue, getNameForChannel(queue), false, false, false, false, nil)

    for d := range msgs {
        // processMessage using the context
        processMessage(ctx, d)
    }
}


```
