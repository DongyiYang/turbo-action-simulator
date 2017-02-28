package api

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type Action struct {
	ProbeType        string                `json:"probeType,omitempty"`
	Account          []*proto.AccountValue `json:"account,omitempty"`
	ActionType       string                `json:"type,omitempty"`
	TargetEntityType string                `json:"entityType,omitempty"`
	TargetEntityID   string                `json:"entityID,omitempty"`
	MoveSpec         *MoveSpec             `json:"moveSpec,omitempty"`
}

//type AccountValue struct {
//	AccountName string `json:"accountName,omitempty"`
//	//// Set of property value lists
//	//GroupScopePropertyValues []*AccountValue_PropertyValueList `protobuf:"bytes,3,rep,name=groupScopePropertyValues" json:"groupScopePropertyValues,omitempty"`
//}

type MoveSpec struct {
	DestinationEntityType string `destinationEntityType,omitempty`
	DestinationEntityID   string `destinationEntityID,omitempty`
	MoveDestinationIP     string `json:"destinationIP,omitempty"`
}
