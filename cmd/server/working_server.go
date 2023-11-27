package main

import (
	server "edetector_go/internal/server"
)

var version string

func main() {
	// go pprof.Pprof_service()
	server.Main(version, nil)
	// stop.Stop()
}
