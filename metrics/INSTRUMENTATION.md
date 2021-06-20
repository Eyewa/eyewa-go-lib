# Instrumentation
Instruments are used to report Measurements. Each Instrument will have the following information:

The **name** of the Instrument \
The **kind** of the Instrument - whether it is a Counter or other instruments, whether it is synchronous or asynchronous \
An optional **unit** of measure\
An optional **description** \
See also [Instrument Naming Convetions](#instrument-naming-convetions)

### Instrumentation Type
Instruments can be categorized based on whether they are synchronous or asynchronous:

**Synchronous instruments** (e.g. **Counter**) are meant to be invoked inline with application/business processing logic. For example, an HTTP client could use a Counter to record the number of bytes it has received. Measurements recorded by synchronous instruments can be associated with the Context.\
**Asynchronous instruments** (e.g. **Asynchronous Gauge**) give the user a way to register callback function, and the callback function will only be invoked upon collection. For example, a piece of embedded software could use an asynchronous gauge to collect the temperature from a sensor every 15 seconds, which means the callback function will only be invoked every 15 seconds. Measurements recorded by asynchronous instruments cannot be associated with the Context.

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
