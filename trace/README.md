# Package trace

## Getting Started

Set the environmental variable `TRACE_EXPORTER_ENDPOINT` if you would like to export to the open telemetry collector, if not present it will default print traces to stdout.
Call `trace.ConfigureAndConnect()` to configure and connect to the collector. Thereafter you would need to add the grpc interceptor to your grpc server.

## GRPC Server

Apply tracing to a server by adding the following server options.

### RPC Server - Request Response

```go
    s := grpc.NewServer(
        trace.UnaryServerTraceInterceptor()
    )
```

### RPC Server - Stream

```go
    s := grpc.NewServer(
        trace.StreamServerTraceInterceptor()
    )
```

## GRPC Client - Request Response

Apply tracing to a client by adding the following dial options.

### RPC Client - Request Response

```go
    conn, err := grpc.Dial(":7777",
        grpc.WithInsecure(),
        trace.UnaryClientTraceInterceptor(),
    )
```

### RPC Client - Stream

```go
    conn, err := grpc.Dial(":7777",
        grpc.WithInsecure(),
        trace.StreamClientTraceInterceptor(),
    )
```
