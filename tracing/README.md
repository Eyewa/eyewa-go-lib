# eyewa-go-lib

Shared Go Lib for Eyewa's microservices.

## tracing

This package configures Open Telemetry as the global tracing provider. It configures an endpoint to where all traces end up using the `TRACE_COLLECTOR_ENDPOINT` env variable. A user of this package is able to view traces on Grafana Tempo.

</br>

### Note

Due to the Open Telemetry trace API and SDK currently at v1.0.0-rc.1, there will be future breaking changes which would require refactoring of trace pkg internals.

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

Env Variable      |Required|
|------------------|------------------|--------|
|SERVICE_NAME                       |y       |Name of the service/application.
|SERVICE_VERSION                    |n       |Version of the service/application.
|EXPORTER_ENDPOINT                  |n       |The endpoint tracing spans are exported to.
|EXPORTER_BLOCKING                  |n       |Exporter initiates a blocking request to an endpoint.
|EXPORTER_SECURE                    |n       |Exporter connects with TLS secure connection.

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
