# eyewa-go-lib

Shared Go Lib for Eyewa's microservices.

## tracing

This package configures Open Telemetry as the global tracing provider. It configures an endpoint to where all traces end up using the `EXPORTER_ENDPOINT` env variable. A user of this package is able to view traces on Grafana Tempo once a trace has been exported.

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
 os.Setenv("EXPORTER_ENDPOINT", "open-telemetry.collector.endpoint")
 os.Setenv("GRPC_SERVER_PORT", "7777")

 err := config.Init()
 if err != nil {
  log.Error(err.Error())
 }

 // launch tracing to open a connection to
 // a tracing backend.
 shutdown, err := tracing.Launch()
 if err != nil {
  log.Error(err.Error())
 }
 defer shutdown()

 // listen on the grpc server port.
 port := os.Getenv("GRPC_SERVER_PORT")
 lis, err := net.Listen("tcp", port)
 if err != nil {
  log.Fatal(err.Error())
 }

 // inject tracing interceptors.
 s := grpc.NewServer(
  tracing.UnaryServerTraceInterceptor(),
  tracing.StreamServerTraceInterceptor(),
 )

 // register the server and start serving grpc requests.
 api.RegisterHelloServiceServer(s, &server{})
 if err := s.Serve(lis); err != nil {
  log.Fatal(err.Error())
 }

}

```

</br>
