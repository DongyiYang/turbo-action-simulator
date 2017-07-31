package api

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type TypeMeta struct {
	Kind string `json:"kind,omitempty"`
}

type Action struct {
	TypeMeta         `json:",inline"`
	ProbeType        string                `json:"probeType,omitempty"`
	Account          []*proto.AccountValue `json:"account,omitempty"`
	ActionType       string                `json:"type,omitempty"`
	TargetEntityType string                `json:"entityType,omitempty"`
	TargetEntityID   string                `json:"entityID,omitempty"`
	MoveSpec         MoveSpec              `json:"moveSpec,omitempty"`
	ScaleSpec        ScaleSpec             `json:"scaleSpec,omitempty"`
}

//type AccountValue struct {
//	AccountName string `json:"accountName,omitempty"`
//	//// Set of property value lists
//	//GroupScopePropertyValues []*AccountValue_PropertyValueList `protobuf:"bytes,3,rep,name=groupScopePropertyValues" json:"groupScopePropertyValues,omitempty"`
//}

type MoveSpec struct {
	DestinationEntityType string `json:"destinationEntityType,omitempty"`
	DestinationEntityID   string `josn:"destinationEntityID,omitempty"`
	MoveDestinationIP     string `json:"destinationIP,omitempty"`
}

type ScaleSpec struct {
	//Provider ProviderInfo `json:"provider,omitempty"`
	ProviderEntityType string `json:"providerEntityType,omitempty"`
	ProviderEntityID   string `json:"providerEntityID,omitempty"`
}

type ProviderInfo struct {
	ProviderEntityType string `json:"providerEntityType,omitempty"`
	ProviderEntityID   string `json:"providerEntityID,omitempty"`
}

type Discovery struct {
	ProbeType     string               `json:"probeType,omitempty"`
	AccountValues []proto.AccountValue `json:"accountValues,omitempty"`
}