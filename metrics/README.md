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

Create a Prometheus exporter.
```go
package demo

import (
	"github.com/eyewa/eyewa-go-lib/log"
	"github.com/eyewa/eyewa-go-lib/metrics"
	"github.com/eyewa/eyewa-go-lib/metrics/prometheus"
	"time"
)

func main() {
	option := prometheus.ExportOption{
		CollectPeriod: 1 * time.Second,
	}
	
	exporter, err := prometheus.NewPrometheusExporter(option)
	if err != nil {
		log.Error(metrics.FailedToInitPrometheusExporterError.Inner(err).Error())
	}
}
```
Create metric launcher with predefined Exporter
```go
ml := metrics.NewMetricLauncher(exporter)
```
Set Global Meter Provider. It will set the Exporter's Meter Provider globally. See also [Global Setting](#global-setting).
```go
ml.SetMeterProvider()
```
Enable Host Instrumentation. See also [Host Instrumentation Metrics](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/host@v0.20.0#pkg-overview) 
```go
ml.EnableHostInstrumentation()
```
Enable Runtime Instrumentation. See also [Runtime Instrumentation Metrics](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/runtime@v0.20.0#pkg-overview)
```go
ml.EnableRuntimeInstrumentation()
```
Then Launch. Launch will start Metric Server on port 2222 on different goroutine \
not to block main process, so it returns error channel to check that is everything \
alright. It is receive-only channel.
```go
errCh := ml.Launch()
```
or put them together
```go
 metrics.NewMetricLauncher(exporter).
    SetMeterProvider().
    EnableHostInstrumentation().
    EnableRuntimeInstrumentation().
    Launch()
	
```
# Instrumentation
Please see [Instrumentation](INSTRUMENTATION.md) from here

---
### Global Setting

When using OpenTelemetry, it’s a good practice to set a global tracer provider and 
a global meter provider. Doing so will make it easier for libraries and other dependencies 
that use the OpenTelemetry API to easily discover the SDK, and emit telemetry data.\
[Setting Global Option](https://opentelemetry.io/docs/go/getting-started/#setting-global-options)

---
### Programming Model under the hood
You can check Open Telemetry programming model here. [Programming Model](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/README.md#programming-model)

---
Read more [Metrics API](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/api.md)
