# eyewa-go-lib

Shared Go Lib for Eyewa's microservices.

## tracing

This package configures open telemetry as the global trace provider. It configures the location of where all traces end up using the `TRACE_COLLECTOR_ENDPOINT`. A user of this package is able to view traces on Grafana Tempo.

</br>

## How To Use

- Set the `TRACE_COLLECTOR_ENDPOINT` environmental variable.
- Instantiate a launcher and launch to connect to the open telemetry collector.
- Add a GRPC interceptor to the GRPC server/client.

</br>

### Configuration Init Order

Upon initialising a tracing configuration, the order in which the configuration is configured is as follows:

1. Environmental Variables
2. Options

</br>

### Environmental Variables

</br>

```go
// the name of the microservice.
"SERVICE_NAME",
// the kubernetes pod/host.
"HOST_NAME",
// the endpoint of the collector traces get exported to.
"TRACE_COLLECTOR_ENDPOINT",
```

  </br>

### Tracing A GRPC Server

```go
launcher := tracing.NewLauncher()
err := launcher.Launch()
if err != nil {
    log.Fatal("failed to launch")
}
defer launcher.Shutdown()

s := grpc.NewServer(
    // trace all unary requests
    trace.UnaryServerTraceInterceptor(),
    // trace all bidirectional streams 
    trace.StreamServerTraceInterceptor(),
)
```

</br>

### Tracing a GRPC client

```go
launcher := tracing.NewLauncher()
err := launcher.Launch()
if err != nil {
    log.Fatal("failed to launch")
}
defer launcher.Shutdown()

port := os.Getenv("GRPC_DIAL_PORT")
conn, err := grpc.Dial(port,
    // trace all unary requests
    grpc.WithInsecure(),
    trace.UnaryClientTraceInterceptor(),
    // trace all bidirectional streams 
    trace.StreamClientTraceInterceptor(),
)
```

</br>
