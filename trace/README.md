# eyewa-go-lib

Shared Go Lib for Eyewa's microservices.

# tracing

This package configures open telemetry as the global trace provider. It configures the location of where all traces end up using the `TRACE_COLLECTOR_ENDPOINT`. The user of this package is able to view traces on Grafana Tempo.

# How to use

The following variables should be injected.

```go
        // the name of the microservice.
        "SERVICE_NAME",
        // the kubernetes pod/host.
        "HOST_NAME",
        // the endpoint of the collector traces get exported to.
        "TRACE_COLLECTOR_ENDPOINT",
```

Set the environmental variable `TRACE_COLLECTOR_ENDPOINT` when exporting to the open telemetry collector else it will default to printing spans to stdout.
Call `trace.ConfigureAndConnect()` to configure and connect to the collector. Thereafter adding a GRPC interceptor to your grpc server.

```go
    trace.ConfigureAndConnect()

    //...

    // apply tracing to server interceptors
    s := grpc.NewServer(
        trace.UnaryServerTraceInterceptor(),
        trace.StreamServerTraceInterceptor(),
    )
```

```go
    trace.ConfigureAndConnect()

    // apply tracing to client interceptors
    conn, err := grpc.Dial(":7777",
        grpc.WithInsecure(),
        trace.UnaryClientTraceInterceptor(),
        trace.StreamClientTraceInterceptor(),
    )
```
