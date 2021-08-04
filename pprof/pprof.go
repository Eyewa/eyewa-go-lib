package pprof

import (
	"net/http"
	_ "net/http/pprof"
)

func init() {
	// initiate pprof
	go func() {
		_ = http.ListenAndServe(":9111", nil)
	}()
}
