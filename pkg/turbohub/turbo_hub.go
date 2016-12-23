package turbohub

import (
	"fmt"

	"github.com/vmturbo/vmturbo-go-sdk/pkg/proto"

	"github.com/turbonomic/turbo-simulator/pkg/mediationcontainer"
	"github.com/turbonomic/turbo-simulator/pkg/rest"

	"github.com/golang/glog"
	goproto "github.com/golang/protobuf/proto"
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
			case webSocketReceivedMessage := <-h.mediationContainer.ReceiveMessage():
				err := h.handleRawClientMessage(webSocketReceivedMessage)
				if err != nil {
					glog.Errorf("Error handling received client message: %s", err)
				}
			case restAPIReceivedMessage := <-h.restManager.ReceiveMessage():
				err := h.sendServerMessage(restAPIReceivedMessage)
				if err != nil {
					glog.Errorf("Error forwarding mediation server message from REST API to "+
						"WebSocket: %s", err)
				}
			case <-h.StopChan:
				return
			}
		}
	}()
}

// Marshall the message into byte array and send the message via mediation container.
func (h *TurboHub) sendServerMessage(serverMsg *proto.MediationServerMessage) error {
	glog.V(3).Infof("Send out to WebSocket: %++v", serverMsg)
	rawServerMsg, err := marshallServerMessage(serverMsg)
	if err != nil {
		return err
	}
	h.mediationContainer.SendMessage(rawServerMsg)
	return nil
}

// Get raw message from mediation container and unmarshall it into MediationClientMessage.
func (h *TurboHub) handleRawClientMessage(rawMessage []byte) error {
	clientMessage, err := unmarshallClientMessage(rawMessage)
	if err != nil {
		return err
	}
	h.forwardClientMessage(clientMessage)
	return nil
}

// Forward message to different component based on message type.
func (h *TurboHub) forwardClientMessage(clientMsg *proto.MediationClientMessage) {
	glog.V(3).Infof("Get client message: %++v", clientMsg)
	if clientMsg.ValidationResponse != nil {
		// TODO
	} else if clientMsg.DiscoveryResponse != nil {
		// TODO
	} else if clientMsg.KeepAlive != nil {
		// TODO
	} else if clientMsg.ActionProgress != nil {
		// TODO
	} else if clientMsg.ActionResponse != nil {
		// TODO
	}

}

func unmarshallClientMessage(rawMessage []byte) (*proto.MediationClientMessage, error) {
	clientMessage := &proto.MediationClientMessage{}
	err := goproto.Unmarshal(rawMessage, clientMessage)
	if err != nil {
		return nil, fmt.Errorf("Cannot unmarshall: %s", err)
	}

	return clientMessage, nil
}

func marshallServerMessage(serverMessage *proto.MediationServerMessage) ([]byte, error) {
	marshalled, err := goproto.Marshal(serverMessage)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling server message %+v", serverMessage)
	}
	return marshalled, nil
}
