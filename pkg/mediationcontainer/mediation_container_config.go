package mediationcontainer


type MediationContainerConfig struct{
	SendMessageChan chan []byte
	ReceiveMessageChan chan []byte

	StopChan chan struct{}
}

func NewMediationContainerConfig() *MediationContainerConfig {
	return &MediationContainerConfig{
		SendMessageChan: make(chan []byte),
		ReceiveMessageChan: make(chan []byte),

		StopChan: make(chan struct{}),
	}
}