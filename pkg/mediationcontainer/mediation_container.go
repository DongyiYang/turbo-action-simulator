package mediationcontainer

import (
	"errors"
	"fmt"

	"github.com/turbonomic/turbo-go-sdk/pkg/proto"

	"github.com/golang/glog"
	goproto "github.com/golang/protobuf/proto"
	"github.com/turbonomic/turbo-go-sdk/pkg/version"
	"golang.org/x/net/websocket"
)

const (
	StatusReady   MediationContainerStatus = "Ready"
	StatusWaiting MediationContainerStatus = "Waiting Client"
	StatusClosed  MediationContainerStatus = "Closed"
)

type MediationContainerStatus string

type MediationContainer struct {
	wsConn *websocket.Conn

	config *MediationContainerConfig

	Status MediationContainerStatus

	pipeline *Pipeline
}

func NewMediationContainer(config *MediationContainerConfig) *MediationContainer {
	pipeline := NewPipeline()
	pipeline.Push(&negotiationMessageHandler{})
	pipeline.Push(&registrationMessageHandler{})
	pipeline.Push(&mediationClientMessageHandler{})

	return &MediationContainer{
		config: config,

		Status: StatusWaiting,

		pipeline: pipeline,
	}
}

func (mc *MediationContainer) OnWebSocketConnected(ws *websocket.Conn) {
	mc.wsConn = ws
	go mc.listenSend()
	go mc.listenReceive()
	go mc.handleReceivedMessage()
	mc.Status = StatusReady
	select {}

}

// Listening any message from client.
func (mc *MediationContainer) listenReceive() {
	glog.V(3).Info("Listening message from client...")
	var err error
	for {
		select {
		case <-mc.config.StopChan:
			return
		default:
		}

		var requestContent []byte
		if err = websocket.Message.Receive(mc.wsConn, &requestContent); err != nil {
			// If WebSocket connection get disconnected, stop the for-loop.
			glog.Errorf("Error receive message: %s", err)
			glog.Error("Client disconnected..")
			return
		}
		mc.config.ReceiveMessageChan <- requestContent
	}
}

// Listening if there is any message server wants to send to client.
func (mc *MediationContainer) listenSend() {
	glog.V(3).Info("Listening message sent to client...")
	var err error
	for {
		select {
		case <-mc.config.StopChan:
			mc.Status = StatusClosed
			return
		case replyContent := <-mc.config.SendMessageChan:
			glog.V(4).Infof("got message on SendMessageChan: %v", replyContent)

			if mc.wsConn == nil {
				glog.Error("websocket is not ready.")
			}
			if err = websocket.Message.Send(mc.wsConn, replyContent); err != nil {
				glog.Errorf("Failed to send message via WebSocket: %s", err)
				break
			}
		}
	}
}

func (mc *MediationContainer) sendMessage(message []byte) {
	mc.config.SendMessageChan <- message
}

func (mc *MediationContainer) receiveMessage() <-chan []byte {
	return mc.config.ReceiveMessageChan
}

// Marshall the message into byte array and send the message via mediation container.
func (mc *MediationContainer) SendServerMessage(serverMsg goproto.Message) error {
	if mc.Status != StatusReady {
		return errors.New("Medation container is not ready.")
	}
	glog.V(3).Infof("Send out to WebSocket: %++v", serverMsg)
	rawServerMsg, err := marshallServerMessage(serverMsg)
	if err != nil {
		return err
	}
	mc.sendMessage(rawServerMsg)
	return nil
}

func (mc *MediationContainer) ReceiveMediationClientMessage() <-chan *proto.MediationClientMessage {
	return mc.config.MediationClientMessageChan
}

func (mc *MediationContainer) handleReceivedMessage() {
	for {
		select {
		case <-mc.config.StopChan:
			mc.Status = StatusClosed
			return
		case rawMessage, ok := <-mc.receiveMessage():
			if !ok {
				// TODO status?
				return
			}
			handler, err := mc.pipeline.Peek()
			if err != nil {
				glog.Errorf("Error handling raw client message: %s", err)
				break
			}
			switch handler.(type) {
			case *negotiationMessageHandler:
				negotiationAnswer, err := handler.HandleRawMessage(rawMessage)
				if err != nil {
					glog.Errorf("Negotiation failed: %s", negotiationAnswer)
					break
				}
				err = mc.SendServerMessage(negotiationAnswer)
				if err != nil {
					glog.Errorf("Failed to send negotiation response: %s", err)
				} else {
					// only handle version negotiation once.
					mc.pipeline.Pop()
				}
			case *registrationMessageHandler:
				registrationAck, err := handler.HandleRawMessage(rawMessage)
				if err != nil {
					glog.Errorf("Registration failed: %s", err)
					break
				}
				err = mc.SendServerMessage(registrationAck)
				if err != nil {
					glog.Errorf("Failed to send negotiation response: %s", err)
				} else {
					// only handle registration once.
					mc.pipeline.Pop()
				}
			case *mediationClientMessageHandler:
				// Handle MediationClientMessages.
				mediationClientMessage, err := handler.HandleRawMessage(rawMessage)
				if err != nil {
					glog.Errorf("%s", err)
				} else {
					msg, ok := mediationClientMessage.(*proto.MediationClientMessage)
					if !ok {
						glog.Errorf("Not a mediation client message: %s", err)
					} else {
						mc.config.MediationClientMessageChan <- msg
					}
				}
			}
		}
	}
}

type negotiationMessageHandler struct{}

func (nh *negotiationMessageHandler) HandleRawMessage(rawMessage []byte) (goproto.Message, error) {
	msg, err := unmarshalNegotiationMessage(rawMessage)
	if err != nil {
		return nil, err
	}
	// TODO
	_, ok := msg.(*version.NegotiationRequest)
	if !ok {
		return nil, errors.New("Not a mediation client message")
	}
	acceptEverything := version.NegotiationAnswer_ACCEPTED
	description := "Accept everything."
	return &version.NegotiationAnswer{
		NegotiationResult: &acceptEverything,
		Description:       &description,
	}, nil
}

func unmarshalNegotiationMessage(rawMessage []byte) (goproto.Message, error) {
	clientMessage := &version.NegotiationRequest{}
	err := goproto.Unmarshal(rawMessage, clientMessage)
	if err != nil {
		return nil, fmt.Errorf("Cannot unmarshall: %s", err)
	}

	return clientMessage, nil
}

type registrationMessageHandler struct{}

func (nh *registrationMessageHandler) HandleRawMessage(rawMessage []byte) (goproto.Message, error) {
	msg, err := unmarshalRegistrationMessage(rawMessage)
	if err != nil {
		return nil, err
	}
	_, ok := msg.(*proto.ContainerInfo)
	if !ok {
		return nil, errors.New("Not a registration message")
	}
	// Ack everything.
	return &proto.Ack{}, nil
}

func unmarshalRegistrationMessage(rawMessage []byte) (goproto.Message, error) {
	clientMessage := &proto.ContainerInfo{}
	err := goproto.Unmarshal(rawMessage, clientMessage)
	if err != nil {
		return nil, fmt.Errorf("Cannot unmarshall: %s", err)
	}

	return clientMessage, nil
}

type mediationClientMessageHandler struct{}

func (mch *mediationClientMessageHandler) HandleRawMessage(rawMessage []byte) (goproto.Message, error) {
	msg, err := unmarshalMediationClientMessage(rawMessage)
	if err != nil {
		return nil, err
	}
	clientMessage, ok := msg.(*proto.MediationClientMessage)
	if !ok {
		return nil, errors.New("Not a mediation client message")
	}
	return clientMessage, nil
}

func unmarshalMediationClientMessage(rawMessage []byte) (goproto.Message, error) {
	clientMessage := &proto.MediationClientMessage{}
	err := goproto.Unmarshal(rawMessage, clientMessage)
	if err != nil {
		return nil, fmt.Errorf("Cannot unmarshall: %s", err)
	}

	return clientMessage, nil
}

func marshallServerMessage(serverMessage goproto.Message) ([]byte, error) {
	marshaled, err := goproto.Marshal(serverMessage)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling server message %+v", serverMessage)
	}
	return marshaled, nil
}
