package mediationcontainer

import "github.com/turbonomic/turbo-go-sdk/pkg/proto"

type MediationContainerConfig struct {
	SendMessageChan    chan []byte
	ReceiveMessageChan chan []byte

	MediationClientMessageChan chan *proto.MediationClientMessage

	StopChan chan struct{}
}

func NewMediationContainerConfig() *MediationContainerConfig {
	return &MediationContainerConfig{
		SendMessageChan:            make(chan []byte),
		ReceiveMessageChan:         make(chan []byte),
		MediationClientMessageChan: make(chan *proto.MediationClientMessage),

		StopChan: make(chan struct{}),
	}
}
