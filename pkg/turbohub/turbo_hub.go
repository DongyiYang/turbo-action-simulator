package turbohub

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/turbonomic/turbo-simulator/pkg/converter"
	"github.com/turbonomic/turbo-simulator/pkg/mediationcontainer"
	"github.com/turbonomic/turbo-simulator/pkg/rest"
	"github.com/turbonomic/turbo-simulator/pkg/rest/api"

	"github.com/turbonomic/turbo-go-sdk/pkg/proto"

	"github.com/golang/glog"
)

type TurboHub struct {
	mediationContainer *mediationcontainer.MediationContainer
	restManager        *rest.RESTManager
	StopChan           chan struct{}
}

func NewTurboHub(mc *mediationcontainer.MediationContainer, rm *rest.RESTManager) *TurboHub {
	return &TurboHub{
		mediationContainer: mc,
		restManager:        rm,
		StopChan:           make(chan struct{}),
	}
}

func (h *TurboHub) Run() {
	go func() {
		for {
			select {
			case webSocketReceivedMessage := <-h.mediationContainer.ReceiveMediationClientMessage():
				err := h.forwardClientMessage(webSocketReceivedMessage)
				if err != nil {
					glog.Errorf("Error handling received client message: %s", err)
				}
			case apiObjectReceived := <-h.restManager.ReceiveMessage():
				err := h.distributeAPIRequest(apiObjectReceived)
				if err != nil {
					glog.Errorf("Error: %s", err)
				}
			case <-h.StopChan:
				return
			}
		}
	}()
}

// Find the type of each APIObject, then distribute them accordingly.
func (h *TurboHub) distributeAPIRequest(apiObj api.APIObject) error {
	switch apiObj.(type) {
	case api.Action:
		action := apiObj.(api.Action)
		serverMessage, err := converter.TransformActionRequest(&action)
		if err != nil {
			return fmt.Errorf("Failed to create mediation server message based on given action request: %s",
				err)
		}
		glog.V(4).Infof("Action request is generated: %++v", serverMessage)
		err = h.sendServerMessage(serverMessage)
		if err != nil {
			return fmt.Errorf("Failed to forward mediation server message from REST API to "+
				"WebSocket: %s", err)
		}
	default:
		glog.Errorf("API object type %s is not supported", reflect.TypeOf(apiObj))
	}
	return nil
}

// Send the message via mediation container.
func (h *TurboHub) sendServerMessage(serverMsg *proto.MediationServerMessage) error {
	if h.mediationContainer == nil {
		return errors.New("Medation container is not set")
	}
	err := h.mediationContainer.SendServerMessage(serverMsg)
	return err
}

// Forward message to different component based on message type.
func (h *TurboHub) forwardClientMessage(clientMsg *proto.MediationClientMessage) error {
	glog.V(3).Infof("Get client message: %++v", clientMsg)
	switch clientMsg.MediationClientMessage.(type) {
	case *proto.MediationClientMessage_ValidationResponse: // TODO
	case *proto.MediationClientMessage_DiscoveryResponse: // TODO
	case *proto.MediationClientMessage_KeepAlive: // TODO
	case *proto.MediationClientMessage_ActionResponse: // TODO
	case *proto.MediationClientMessage_ActionProgress: // TODO
	}
	return nil
}
