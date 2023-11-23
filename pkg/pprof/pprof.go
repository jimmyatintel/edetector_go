package pprof

import (
	"net/http"
	_ "net/http/pprof"
)

func Pprof_service() {
    http.ListenAndServe("0.0.0.0:3030", nil)
}