# eyewa-go-lib
Shared Go Lib for Eyewa's microservices.

# metrics
This package is simply a wrapper for OpenTelemetry metric package. It uses Prometheus
Exporter to export metrics. Metrics is exported from **2222** port

# How to use it
The Metrics package consists of the following:
- A Metrics Launcher - serves metrics over HTTP and works with the Metrics Exporter.
- A Metrics Exporter - is used for scrapping data for Prometheus
- A Metrics Instrumentation - any instrumentation of choice to create metrics.

## How to create a metric launcher

```go
package demo

import (
	"context"
	"github.com/eyewa/eyewa-go-lib/errors"
	"github.com/eyewa/eyewa-go-lib/log"
	"github.com/eyewa/eyewa-go-lib/metrics"
	"go.opentelemetry.io/otel/metric"
	"time"
)

func main() {
	option := metrics.ExportOption{
		CollectPeriod: 1 * time.Second,
	}

	ml, err := metrics.NewLauncher(option)
	if err != nil {
		log.Error(errors.FailedToStartMetricServerError.Error())
	}

	ml.SetMeterProvider().
		EnableHostInstrument().
		EnableRuntimeInstrument().
		Launch()

	//Start to create meters
	httpMeter := metrics.NewMeter("http.meter", nil)

	// Create a new instrument from meter
	requestCounter, err := httpMeter.NewCounter("request.counter")
	if err != nil {
		log.Error(errors.FailedToCreateInstrumentError.Error())
	}
	
	// increase measurement
	requestCounter.Add(1)

	// Create async counter with callback
	cb := func(ctx context.Context, result metric.Float64ObserverResult) {
		// increase measurement
		result.Observe(1)
	}

	_, err = httpMeter.NewAsyncCounter("request.async.counter", cb)
	if err != nil {
		log.Error(errors.FailedToCreateInstrumentError.Error())
	}
}
```
It will set the Exporter's Meter Provider globally. See also [Setting Global Option](https://opentelemetry.io/docs/go/getting-started/#setting-global-options)
```go
ml.SetMeterProvider()
```
Enable host instrumentation. See also [Host Instrumentation Metrics](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/host@v0.20.0#pkg-overview) 
```go
ml.EnableHostInstrumentation()
```
Enable runtime instrumentation. See also [Runtime Instrumentation Metrics](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/runtime@v0.20.0#pkg-overview)
```go
ml.EnableRuntimeInstrumentation()
```
Launch will start Metric Server on port 2222 on different goroutine 
not to block main process.
```go
ml.Launch()
```

# Instrumentation
There are two type instrument sync and async. Sync instruments are;
- Counter
- UpDownCounter
- ValueRecorder

asyncs are;
- AsyncCounter
- AsyncUpDownCounter
- AsyncValueRecorder

---
### Programming Model under the hood
Open Telemetry [Programming Model](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/README.md#programming-model)

---
Read more [Metrics API](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/api.md)
