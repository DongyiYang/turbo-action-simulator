package mediationcontainer

import (
	"github.com/golang/glog"
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
}

func NewMediationContainer(config *MediationContainerConfig) *MediationContainer {
	return &MediationContainer{
		config: config,

		Status: StatusWaiting,
	}
}

func (mc *MediationContainer) OnWebSocketConnected(ws *websocket.Conn) {
	mc.wsConn = ws
	go mc.listenSend()
	go mc.listenReceive()
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
			glog.Infof("got message: %v", replyContent)

			if mc.wsConn == nil {
				glog.Error("websocket is not ready.")
			}
			if err = websocket.Message.Send(mc.wsConn, replyContent); err != nil {
				glog.Errorf("Failed to send message via WebSocket: %s", err)
				return
			}
		default:
		}
	}
}

func (mc *MediationContainer) SendMessage(message []byte) {
	mc.config.SendMessageChan <- message
}

func (mc *MediationContainer) ReceiveMessage() <-chan []byte {
	return mc.config.ReceiveMessageChan
}
