# eyewa-go-lib
Shared Go Lib for Eyewa's microservices.

# log
This package provides an abstraction layer for Uber's Zap logger pkg under the hood. For any client requiring logging, it is as simple as setting the `LOG_LEVEL` env var, and then calling the `SetLogLevel` func to initiate the logger.

For each log level supported, there are equivalent log funcs for tracing capabilities if and when required.

```go
package demo

"github.com/eyewa/eyewa-go-lib/log"

func main() {
  os.Setenv("LOG_LEVEL", "debug") // this should be injected and not hardcoded.
  log.SetLogLevel()

  log.Debug("testing")
  log.InfoWithTraceID(uuid.NewString(), "testing 123")
}
```

