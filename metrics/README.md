# eyewa-go-lib
Shared Go Lib for Eyewa's microservices.

# metrics
This package is simply a wrapper for OpenTelemetry metric package.

# How to use it
Basically launch metric server and create new meter from it, then create
new instruments from the meter. After launching, prometheus exporter take
care exporting.

```Go
package demo

import (
	"context"
	"github.com/eyewa/eyewa-go-lib/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"math/rand"
	"sync"
	"time"
)

func main() {
	_ := metrics.Launch(metrics.Option{
		CollectPeriod: 100 * time.Millisecond,
	})

	meter := metrics.NewMeter("custom.meter", nil)

	counter := meter.NewCounter("customer.counter", attribute.Any("version", "1.0.0"))
	counter.Add(1, attribute.Any("label", "new request arrived"))

	updownCounter := meter.NewUpDownCounter("ilk_updown_counter",
		attribute.Any("description", "new updown counter"))
	updownCounter.Add(1)

	valueRecorder := meter.NewValueRecorder("ilk_value_recorder")
	valueRecorder.Record(14)

	observerLock := new(sync.RWMutex)
	observerValueToReport := new(float64)
	observerLabelsToReport := new([]attribute.KeyValue)
	callback := func(ctx context.Context, result metric.Float64ObserverResult) {
		(*observerLock).RLock()
		value := *observerValueToReport
		labels := *observerLabelsToReport
		(*observerLock).RUnlock()
		result.Observe(value, labels...)
	}

	_ = meter.NewSumObserver("sum_observer", callback)
	_ = meter.NewUpDownSumObserver("updownsum_observer", callback)
	_ = meter.NewValueObserver("value_observer", callback)
}
```
---
// startHostInstrument starts Host instrumentation.
// https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/host@v0.20.0
---
// startHostInstrument starts Runtime instrumentation.
// https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/runtime@v0.20.0
---
### Programming Model
https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/README.md#programming-model

---
Additive instruments may be monotonic, in which case they are non-decreasing and naturally define a rate.

The synchronous instrument names are: 

**Counter**:           additive, monotonic \
**UpDownCounter**:     additive \
**ValueRecorder**:     grouping 

and the asynchronous instruments are: 

**SumObserver**:       additive, monotonic \
**UpDownSumObserver**: additive \
**ValueObserver**:     grouping \
https://pkg.go.dev/go.opentelemetry.io/otel/metric

---
### Instrument Naming Convetions
Instruments are associated with the Meter during creation, and are identified by the name:

Meter implementations MUST return an error when multiple Instruments are registered under 
the same Meter instance using the same name.Different Meters MUST be treated as separate namespaces.
The names of the Instruments under one Meter SHOULD NOT interfere with Instruments under another Meter.

Instrument names MUST conform to the following syntax (described using the Augmented Backus-Naur Form):

**instrument-name** = **ALPHA** 0*62 ("_" / "." / "-" / **ALPHA** / **DIGIT**)

**ALPHA** = %x41-5A / %x61-7A; A-Z / a-z \
**DIGIT** = %x30-39 ; 0-9 \
They are not null or empty strings. \
They are case-insensitive, ASCII strings. \
The first character must be an alphabetic character. \
Subsequent characters must belong to the alphanumeric characters, '_', '.', and '-'. \
They can have a maximum length of 63 characters.

---
### Global setting
When using OpenTelemetry, it’s a good practice to set a global tracer provider and 
a global meter provider. Doing so will make it easier for libraries and other dependencies 
that use the OpenTelemetry API to easily discover the SDK, and emit telemetry data.\
https://opentelemetry.io/docs/go/getting-started/

