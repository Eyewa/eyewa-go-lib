# eyewa-go-lib

Shared Go Lib for Eyewa's microservices.

## Docs

- [How tracing works](https://eyewadxb.atlassian.net/wiki/spaces/TECH/pages/1978564609/How+Tracing+Works)
- [Tracing within the RabbitMQ Client](https://eyewadxb.atlassian.net/wiki/spaces/TECH/pages/1880981528/RabbitMQ+Tracing)

## tracing

This package enables Open Telemetry as the global tracing provider. It provides an endpoint where all services using the 'EXPORTER_ENDPOINT' env variable can send their trace payloads to. A user of this package is able to view traces on Grafana Tempo once a trace has been exported.

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
HOSTNAME // The hostname of the deployment instance. | #Optional

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
