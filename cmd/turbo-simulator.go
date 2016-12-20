package main

import(
	"flag"

	"github.com/turbonomic/turbo-action-simulator/cmd/server"

)


func init() {

	flag.Set("logtostderr", "true")
}

func main() {
	flag.Parse()

	ss := &server.SimulatorServer{}
	ss.Run()
}
