package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"

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

	addProfiler(router)

	// Register routes and handle function for http REST request.
	for _, route := range rest.RoutesPaths {
		router.HandleFunc(route.Path, restManager.HandleRequest)
	}

	glog.V(2).Info("Turbo simulator is started.")

	if err := http.ListenAndServe(":1234", router); err != nil {
		log.Fatalf("ListenAndServe with error: %s", err)
	}
	select {}
}

// Add profiler routes to router.
func addProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	// Manually add support for paths linked to by index page at /debug/pprof/
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
}
