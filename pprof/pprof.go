package pprof

import (
	"net/http"
	_ "net/http/pprof"
	"sync"
)

func init() {
	new(sync.Once).Do(func() {
		// initiate pprof
		go func() {
			_ = http.ListenAndServe(":9111", nil)
		}()
	})
}
