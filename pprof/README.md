# eyewa-go-lib
Shared Go Lib for Eyewa's microservices.

# pprof
This package exposes Go's profiling capability via its HTTP server runtime profiling data in the format expected by the pprof visualization tool. This 

Default port exposed is `9111`. 

# How to use
In the **main** goroutine of a service/application the `pprof` package can be imported for its side effects.

```go
package main

import (
	...
	_ "github.com/eyewa/eyewa-go-lib/pprof" // pprof
	...
)
```

This can now reachable on: `http://locahost:9111/debug/pprof`.
