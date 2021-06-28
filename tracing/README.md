# eyewa-go-lib

Shared Go Lib for Eyewa's microservices.

## tracing

This package configures Open Telemetry as the global tracing provider. It configures an endpoint to where all traces end up using the `TRACE_COLLECTOR_ENDPOINT` env variable. A user of this package is able to view traces on Grafana Tempo once a trace has been exported.

</br>

### Note

Due to the Open Telemetry trace API and SDK currently at v1.0.0-rc.1, there will be future breaking changes which would require refactoring of trace pkg internals.

</br>

## How To Use

- Set the `EXPORTER_ENDPOINT` environmental variable.
- launch to connect to the open telemetry collector.
- Add a GRPC interceptor to the GRPC server/client.

</br>

### Environmental Variables

```go
SERVICE_NAME // Name of the service/application. #Required
SERVICE_VERSION // Version of the service/application.
EXPORTER_ENDPOINT // The endpoint that spans get exported to.
EXPORTER_BLOCKING // Exporter initiates a blocking request to an endpoint.
EXPORTER_SECURE // Exporter connects with TLS secure connection.
```

</br>

### Tracing A GRPC Server

```go
shutdown, err := tracing.Launch()
if err != nil {
    log.Fatal("failed to launch tracing")
}
defer shutdown()

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
shutdown, err := tracing.Launch()
if err != nil {
    log.Fatal("failed to launch tracing")
}
defer shutdown()

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
