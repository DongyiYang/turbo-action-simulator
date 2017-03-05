package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/turbonomic/turbo-simulator/pkg/mediationcontainer"
	"github.com/turbonomic/turbo-simulator/pkg/rest"
	"github.com/turbonomic/turbo-simulator/pkg/turbohub"

	"github.com/golang/glog"
	"github.com/gorilla/mux"

	"golang.org/x/net/websocket"
)

type SimulatorServer struct{}

func (s *SimulatorServer) Run() {
	config := mediationcontainer.NewMediationContainerConfig()
	mediationContainer := mediationcontainer.NewMediationContainer(config)

	restManager := rest.NewRESTManager()

	turboHub := turbohub.NewTurboHub(mediationContainer, restManager)
	turboHub.Run()

	router := mux.NewRouter()
	router.HandleFunc("/vmturbo/remoteMediation", func(w http.ResponseWriter, r *http.Request) {
		glog.V(4).Infof("request is %+v\n", r)
		headers := r.Header
		// TODO maybe there is a better way to check this is a WebSocket connection.
		if _, exist := headers["Sec-Websocket-Key"]; exist {
			glog.V(4).Info("This is a websocket connection.")
			websocket.Handler(mediationContainer.OnWebSocketConnected).ServeHTTP(w, r)
		} else {
			glog.V(4).Info("A http connection.")
		}
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("A http connection.")
	})

	// Register routes and handle function for http REST request.
	for _, route := range rest.RoutesPaths {
		router.HandleFunc(route.Path, restManager.HandleRequest)
		//for _, method := range route.Method {
		//	r.Methods(method)
		//}
	}

	glog.V(2).Info("Turbo simulator is started.")

	if err := http.ListenAndServe(":1234", router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	select {}
}
