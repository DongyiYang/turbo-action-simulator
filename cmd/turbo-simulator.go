package main

import (
	"flag"

	"github.com/turbonomic/turbo-simulator/cmd/server"
)

func main() {

	ss := server.NewSimulatorServer()
	ss.AddFlags()
	flag.Parse()
	ss.Run()
}
