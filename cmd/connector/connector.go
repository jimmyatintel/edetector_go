package main

import (
	"edetector_go/internal/rbconnector"
)

var version string

func main() {
	rbconnector.Start(version)
}
