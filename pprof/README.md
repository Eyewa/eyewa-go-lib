# eyewa-go-lib
Shared Go Lib for Eyewa's microservices.

# pprof
This package exposes Go's profiling capability via its HTTP server for profiling a service/application.

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

This can then be reachable on: `http://127.0.0.1:9111/debug/pprof`.
