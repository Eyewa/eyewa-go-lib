# eyewa-go-lib
Shared Go Lib for Eyewa's microservices.

# metrics
This package is simply a wrapper for OpenTelemetry metric package. It uses Prometheus
Exporter to export metrics. Metrics are exported from **2222** port

# How to use it
The Metrics package consists of the following:
- A Metrics Launcher - serves metrics over HTTP and works with the Metrics Exporter.
- A Metrics Exporter - is used for scrapping data for Prometheus
- A Metrics Instrumentation - any instrumentation of choice to create metrics.

## How to create a metric launcher
The following variable can be injected in order to use this pkg
```
METRICS_COLLECTOR_INTERVAL=20s  // optional - default is 10s if var is not provided
```

```go
package demo

import (
	"github.com/eyewa/eyewa-go-lib/errors"
	"github.com/eyewa/eyewa-go-lib/log"
	"github.com/eyewa/eyewa-go-lib/metrics"
	"go.opentelemetry.io/otel/metric"
)

func main() {
	ml, err := metrics.NewLauncher()
	if err != nil {
		log.Error(errors.ErrorFailedToStartMetricServer.Error())
	}

	ml.SetMeterProvider().
		EnableHostInstrumentation().
		EnableRuntimeInstrumentation().
		Launch()
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
`Launch` will start the Metrics Server as a goroutine on port `2222` in order to avoid blocking the main process.
```go
ml.Launch()
```
## How to create an instrument
```go
    //Start to create meters 
    httpMeter := metrics.NewMeter("http.meter", nil)

    // Create a new instrument from meter
    requestCounter, err := httpMeter.NewCounter("request.counter")
    if err != nil {
        log.Error(errors.ErrorFailedToCreateInstrument.Error())
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
        log.Error(errors.ErrorFailedToCreateInstrument.Error())
    }
```
# Instruments
There are two types of instrument sync and async. Sync instruments are:
- Counter
- UpDownCounter
- ValueRecorder

asyncs are:
- AsyncCounter
- AsyncUpDownCounter
- AsyncValueRecorder

---
# Best practice on adding metrics instrumentation within a microservice
Instruments should be defined under a custom struct. The custom struct should be 
initialized on top of the service. It is better for readibility and tracking which 
metrics are used for the service. 

For example; define custom instruments under a struct
```go
// create a metrics.go
type CatalogConsumerMetrics struct{
	ProductCreatedEventCounter *metrics.Counter
}

// Initialize it on top of service.
func NewCatalogConsumerMetrics() (*CatalogConsumerMetrics, error){
    meter := NewMeter("catalog.consumer",nil)
    
    productCreatedEventCounter, err := meter.NewCounter("product.created.event.counter")
    if err != nil{
    	return nil, err
    }
    
    return &CatalogConsumerMetrics{
        ProductCreatedEventCounter: productCreatedEventCounter,
    }, nil
}
```
```go
// in main.go
func main(){
    metrics, err := NewCatalogConsumerMetrics()
    if err != nil { 
	    log.Error(erros.ErrorFailedToCreateInstrument.Error())
    }
}
```

---
## Metrics WIKI

For detailed information please see [confluence page](https://eyewadxb.atlassian.net/wiki/spaces/TECH/pages/1869545495/Metrics+Package)

---
### Programming Model under the hood
Open Telemetry [Programming Model](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/README.md#programming-model)

---
Read more [Metrics API](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/api.md)
