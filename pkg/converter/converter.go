package converter

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/turbonomic/turbo-simulator/pkg/rest/api"
	"github.com/turbonomic/turbo-simulator/pkg/turbomessage"

	"github.com/golang/glog"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

var (
	// TODO Support more type of entity and action.
	// Map a string to corresponding EntityDTO_EntityType.
	entityTypeConverter map[string]proto.EntityDTO_EntityType = map[string]proto.EntityDTO_EntityType{
		"VirtualMachine": proto.EntityDTO_VIRTUAL_MACHINE,
		"Pod":            proto.EntityDTO_CONTAINER_POD,
		"Application":    proto.EntityDTO_APPLICATION,
	}

	// Map a string to ActionItemDTO_ActionType
	actionTypeConverter map[string]proto.ActionItemDTO_ActionType = map[string]proto.ActionItemDTO_ActionType{
		"move":      proto.ActionItemDTO_MOVE,
		"provision": proto.ActionItemDTO_PROVISION,
	}
)

// Build an MediationServerMessage with ActionRequest from given action information from API.
func TransformActionRequest(actionAPIRequest *api.Action) (*proto.MediationServerMessage, error) {
	// 1. Get the target entity type.
	tSETypeString := actionAPIRequest.TargetEntityType
	glog.V(3).Infof("Target entity type is: %s", tSETypeString)
	if tSETypeString == "" {
		return nil, errors.New("Target entity type is not provided.")
	}
	targetEntityType, exist := entityTypeConverter[tSETypeString]
	if !exist {
		return nil, fmt.Errorf("Target entity type %s is not supported.", actionAPIRequest.TargetEntityType)
	}
	// 2. Build the target entity dto.
	targetEntityDTO, err := builder.NewEntityDTOBuilder(targetEntityType,
		actionAPIRequest.TargetEntityID).Create()
	if err != nil {
		return nil, err
	}
	// 3. Get the action type.
	actionType, exist := actionTypeConverter[actionAPIRequest.ActionType]
	if !exist {
		return nil, fmt.Errorf("Action type %s is not supported.", actionType)
	}
	// 4. Build the action item dto.
	actionItemDTOBuilder := turbomessage.NewActionItemDTOBuilder(actionType).
		TargetSE(targetEntityDTO)
	switch actionType {
	case proto.ActionItemDTO_MOVE:
		newEntityDTO, err := getMovingTarget(actionAPIRequest.MoveSpec)
		if err != nil {
			return nil, err
		}
		actionItemDTOBuilder.NewSE(newEntityDTO)
	case proto.ActionItemDTO_PROVISION:
		providerInfo, err := getScalingProvider(actionAPIRequest.ScaleSpec)
		if err != nil {
			return nil, err
		}
		actionItemDTOBuilder.Provider(providerInfo)
	}
	actionItemDTO, err := actionItemDTOBuilder.Build()
	if err != nil {
		return nil, err
	}
	// 5. Build the action execution dto.
	actionExecutionDTO, err := turbomessage.NewActionExecutionDTOBuilder(actionType).
		ActionItem(actionItemDTO).Build()
	if err != nil {
		return nil, err
	}
	// 6. Build the action request.
	pType := actionAPIRequest.ProbeType
	if pType == "" {
		pType = turbomessage.DefaultProbeType
	}
	accountValue := actionAPIRequest.Account
	if accountValue == nil || len(accountValue) == 0 {
		accountValue = turbomessage.DefaultAccountValues
	}
	actionRequest, err := turbomessage.NewActionRequestBuilder(pType, accountValue, actionExecutionDTO).Build()
	if err != nil {
		return nil, err
	}

	// TODO message ID is a random number in [0, 1000).
	messageID := rand.Int31n(1000)
	serverMessage := turbomessage.NewMediationServerMessageBuilder(messageID).ActionRequest(actionRequest).Build()
	return serverMessage, nil
}

// Build the move action target based on the given move spec.
func getMovingTarget(moveSpec api.MoveSpec) (*proto.EntityDTO, error) {
	glog.V(3).Infof("Move spec is %++v", moveSpec)
	nSETypeString := moveSpec.DestinationEntityType
	if nSETypeString == "" {
		return nil, errors.New("New service entity type is not provide for move action.")
	}
	newEntityType, exist := entityTypeConverter[nSETypeString]
	if !exist {
		return nil, fmt.Errorf("Destination entity type %s is not supported.",
			moveSpec.DestinationEntityType)
	}

	newEntityDTOBuilder := builder.NewEntityDTOBuilder(newEntityType,
		moveSpec.DestinationEntityID)
	switch newEntityType {
	case proto.EntityDTO_VIRTUAL_MACHINE:
		virtualMachineData := &proto.EntityDTO_VirtualMachineData{
			IpAddress: []string{moveSpec.MoveDestinationIP},
		}
		newEntityDTOBuilder.VirtualMachineData(virtualMachineData)
	}

	newEntityDTO, err := newEntityDTOBuilder.Create()
	if err != nil {
		return nil, err
	}
	return newEntityDTO, nil
}

func getScalingProvider(scaleSpec api.ScaleSpec) (*proto.ActionItemDTO_ProviderInfo, error) {
	glog.V(3).Infof("Scale spec is %++v", scaleSpec)
	if scaleSpec.ProviderEntityType == "" || scaleSpec.ProviderEntityID == "" {
		return nil, errors.New("Required provider info is missing or invalid.")
	}

	entityType, exist := entityTypeConverter[scaleSpec.ProviderEntityType]
	if !exist {
		return nil, fmt.Errorf("Provider entity type %s is not provide for scaling action.", entityType)
	}
	return &proto.ActionItemDTO_ProviderInfo{
		EntityType: &entityType,
		Ids:        []string{scaleSpec.ProviderEntityID},
	}, nil
}
