package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/gorilla/mux"

	//"github.com/turbonomic/turbo-go-sdk/pkg/proto"

	//"github.com/turbonomic/turbo-simulator/pkg/converter"
	"github.com/turbonomic/turbo-simulator/pkg/rest/api"
)

// Handle API request, return the related MeditationServerMessage instance.
type TurboHandleFunc func(w http.ResponseWriter, r *http.Request) error

type APIHandler struct {
	handlers               map[string]TurboHandleFunc

	apiObjectGeneratorChan chan api.APIObject
}

func NewAPIHandler(apiObjectGeneratorChan chan api.APIObject) *APIHandler {
	apiHandler := &APIHandler{
		apiObjectGeneratorChan: apiObjectGeneratorChan,
	}
	// register handlers
	handlers := make(map[string]TurboHandleFunc)
	handlers["actions"] = apiHandler.handleActionRequest

	apiHandler.handlers = handlers
	return apiHandler
}

func (h *APIHandler) handleAPIRequest(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.RequestURI, "/")
	fmt.Printf("Path is %s, with length %d\n", paths, len(paths))
	if len(paths) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Please provide full API request url.")
	}
	entityType := paths[2]
	if len(entityType) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Please provide entity type.")
	}
	entityHandlerFunc, exist := h.handlers[entityType]
	if !exist {
		fmt.Fprintf(w, "API type %s is not supported.", entityType)
	}

	err := entityHandlerFunc(w, r)
	if err != nil {
		glog.Errorf("Got error when handle API request: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server error.\n")
	}
}

// Handle action related API requests.
func (h *APIHandler) handleActionRequest(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	switch r.Method {
	case "GET":
		id := vars["id"]
		if id != "" {
			// TODO find msg based on id.
			glog.V(3).Infof("Get action message with id %s", id)
		} else {
			// TODO list all msg.
			glog.V(3).Info("Get all action messages")
		}
		return nil
	case "POST":
		// TODO create msg.
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("Cannot create action from API: %s", err)
		}
		if err := r.Body.Close(); err != nil {
			return fmt.Errorf("Cannot create action from API: %s", err)
		}

		var action api.Action
		if err := json.Unmarshal(body, &action); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				return fmt.Errorf("Cannot encode error message: %s", err)
			}
		}
		glog.V(3).Infof("Created a new action instance from REST API: %++v", action)

		//serverMessage, err := converter.TransformActionRequest(&action)
		//if err != nil {
		//	return fmt.Errorf("Failed to create mediation server message based on given request: %s",
		//		err)
		//}
		//glog.V(3).Infof("Build mediation server message: %+v", serverMessage)

		// Send action instance to channel, which will then be passed ot rest manager.
		h.apiObjectGeneratorChan <- action
		return nil
	case "DELETE":
		// TODO delete msg.
		return nil
	default:
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("Unsupported method %s", r.Method)
	}
}
