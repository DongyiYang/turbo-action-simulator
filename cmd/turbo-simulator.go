package main

import(
	"flag"

	"github.com/turbonomic/turbo-simulator/cmd/server"

)


func init() {

	flag.Set("logtostderr", "true")
}

func main() {
	flag.Parse()

	ss := &server.SimulatorServer{}
	ss.Run()
}
