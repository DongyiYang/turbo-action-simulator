package rest

import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/golang/glog"

	"github.com/vmturbo/vmturbo-go-sdk/pkg/proto"

	"github.com/turbonomic/turbo-simulator/pkg/rest/api"
	"github.com/turbonomic/turbo-simulator/pkg/converter"
)

// Handle API request, return the related MeditationServerMessage instance.
type TurboHandleFunc func(w http.ResponseWriter, r *http.Request) (*proto.MediationServerMessage, error)

type APIHandler struct {
	handlers map[string]TurboHandleFunc
}

func NewAPIHandler() *APIHandler {
	handlers := make(map[string]TurboHandleFunc)
	handlers["actions"] = handleActionRequest

	return &APIHandler{handlers}
}

func (h *APIHandler) handleAPIRequest(w http.ResponseWriter, r *http.Request) *proto.MediationServerMessage {
	paths := strings.Split(r.RequestURI, "/")
	fmt.Printf("Path is %s, with length %d\n", paths, len(paths))
	if len(paths) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Please provide full API request url.")
		return nil
	}
	entityType := paths[2]
	if len(entityType) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Please provide entity type.")
		return nil
	}
	entityHandlerFunc, exist := h.handlers[entityType]
	if !exist {
		fmt.Fprintf(w, "Entity type %s is not supported.", entityType)
		return nil
	}

	msg, err := entityHandlerFunc(w, r)
	if err != nil {
		glog.Errorf("Got error when handle API request: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Server error.")
		return nil
	}
	return msg
}

func handleActionRequest(w http.ResponseWriter, r *http.Request) (*proto.MediationServerMessage, error) {
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
		return nil, nil
	case "POST":
		// TODO create msg.
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("Cannot create action from API: %s", err)
		}
		if err := r.Body.Close(); err != nil {
			return nil, fmt.Errorf("Cannot create action from API: %s", err)
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var action api.Action
		if err := json.Unmarshal(body, &action); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				return nil, fmt.Errorf("Cannot encode error message: %s", err)
			}
		}
		glog.V(3).Infof("Created a new action instance from REST API: %++v", action)

		serverMessage, err:= converter.TransformActionRequest(&action)
		if err != nil {
			return nil, fmt.Errorf("Failed to create mediation server message based on given request: %s",
				err)
		}
		glog.V(3).Infof("Build mediation server message: %+v", serverMessage)
		return serverMessage, err
	case "DELETE":
		// TODO delete msg.
		return nil, nil
	default:
		w.WriteHeader(http.StatusBadRequest)
		return nil, fmt.Errorf("Unsupported method %s", r.Method)
	}
}
