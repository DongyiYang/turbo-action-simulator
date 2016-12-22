package turbomessage

import (
	"github.com/vmturbo/vmturbo-go-sdk/pkg/proto"
)

type MediationServerMessageBuilder struct {
	validationRequest *proto.ValidationRequest
	discoveryRequest  *proto.DiscoveryRequest
	actionRequest     *proto.ActionRequest
	messageID         *int32
}

func NewMediationServerMessageBuilder(messageID int32) *MediationServerMessageBuilder {
	return &MediationServerMessageBuilder{
		messageID: &messageID,
	}
}

func (mb *MediationServerMessageBuilder) Build() *proto.MediationServerMessage {
	return &proto.MediationServerMessage{
		ValidationRequest: mb.validationRequest,
		DiscoveryRequest:  mb.discoveryRequest,
		ActionRequest:     mb.actionRequest,
		MessageID:         mb.messageID,
	}
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
