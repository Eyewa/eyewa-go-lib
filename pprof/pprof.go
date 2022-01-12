package pprof

import (
	"net/http"
	_ "net/http/pprof"
	"os"
)

func Init() {
	if os.Getenv("ENV") == "dev" {
		// initiate pprof
		go func() {
			_ = http.ListenAndServe(":9111", nil)
		}()
	}
}
