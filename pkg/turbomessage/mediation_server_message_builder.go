package turbomessage

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type MediationServerMessageBuilder struct {
	validationRequest     *proto.ValidationRequest
	discoveryRequest      *proto.DiscoveryRequest
	actionRequest         *proto.ActionRequest
	interruptionOperation *int32
	messageID             *int32
}

func NewMediationServerMessageBuilder(messageID int32) *MediationServerMessageBuilder {
	return &MediationServerMessageBuilder{
		messageID: &messageID,
	}
}

func (mb *MediationServerMessageBuilder) Build() *proto.MediationServerMessage {
	serverMessage := &proto.MediationServerMessage{
		MessageID: mb.messageID,
	}

	if mb.validationRequest != nil {
		serverMessage.MediationServerMessage = &proto.MediationServerMessage_ValidationRequest{
			ValidationRequest: mb.validationRequest,
		}
	} else if mb.discoveryRequest != nil {
		serverMessage.MediationServerMessage = &proto.MediationServerMessage_DiscoveryRequest{
			DiscoveryRequest: mb.discoveryRequest,
		}
	} else if mb.actionRequest != nil {
		serverMessage.MediationServerMessage = &proto.MediationServerMessage_ActionRequest{
			ActionRequest: mb.actionRequest,
		}
	} else if mb.interruptionOperation != nil {
		serverMessage.MediationServerMessage = &proto.MediationServerMessage_InterruptOperation{
			InterruptOperation: *mb.interruptionOperation,
		}
	}

	return serverMessage
}

func (mb *MediationServerMessageBuilder) ActionRequest(
	actionRequest *proto.ActionRequest) *MediationServerMessageBuilder {
	mb.actionRequest = actionRequest
	return mb
}

func (mb *MediationServerMessageBuilder) ValidationRequest(
	validationRequest *proto.ValidationRequest) *MediationServerMessageBuilder {
	mb.validationRequest = validationRequest
	return mb
}

func (mb *MediationServerMessageBuilder) DiscoveryRequest(
	discoveryRequest *proto.DiscoveryRequest) *MediationServerMessageBuilder {
	mb.discoveryRequest = discoveryRequest
	return mb
}

func (mb *MediationServerMessageBuilder) InterruptionOperation(
	operation int32) *MediationServerMessageBuilder {
	mb.interruptionOperation = &operation
	return mb
}
