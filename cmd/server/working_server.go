package main

import (
	server "edetector_go/internal/server"
	pprof_test "edetector_go/pkg/pprof"
)

var version string

func main() {
	go pprof_test.Pprof_service()
	server.Main(version, nil)
	// stop.Stop()
}
