package main

import (
	server "edetector_go/internal/server"
)

var version string

func main() {
	// go pprof_test.Pprof_service()
	server.Main(version, nil)
	// stop.Stop()
}
