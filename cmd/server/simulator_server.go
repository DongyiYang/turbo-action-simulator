package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/turbonomic/turbo-action-simulator/pkg/mediationcontainer"

	"golang.org/x/net/websocket"
	"github.com/turbonomic/turbo-action-simulator/pkg/turbohub"
	"github.com/golang/glog"
)

type SimulatorServer struct{}

func (s *SimulatorServer) Run() {
	config := mediationcontainer.NewMediationContainerConfig()
	mediationContainer := mediationcontainer.NewMediationContainer(config)

	turboHub := turbohub.NewTurboHub(mediationContainer)
	turboHub.Run()

	http.HandleFunc("/vmturbo/remoteMediation", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("request is %+v\n", r)
		headers := r.Header
		// TODO maybe there is a better way to check this is a WebSocket connection.
		if _, exist := headers["Sec-Websocket-Key"]; exist {
			fmt.Println("This is a websocket connection.")
			websocket.Handler(mediationContainer.OnWebSocketConnected).ServeHTTP(w, r)
		} else {
			fmt.Println("A http connection.")
		}
	})

	glog.V(2).Info("Turbo simulator is started.")

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	select {}
}
