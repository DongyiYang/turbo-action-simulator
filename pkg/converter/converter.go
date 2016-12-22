package converter

import (
	"fmt"
	"math/rand"

	"github.com/turbonomic/turbo-action-simulator/pkg/rest/api"
	"github.com/turbonomic/turbo-action-simulator/pkg/turbomessage"

	"github.com/golang/glog"
	"github.com/vmturbo/vmturbo-go-sdk/pkg/builder"
	"github.com/vmturbo/vmturbo-go-sdk/pkg/proto"
)

var (
	entityTypeConverter map[string]proto.EntityDTO_EntityType = map[string]proto.EntityDTO_EntityType{
		"VirtualMachine": proto.EntityDTO_VIRTUAL_MACHINE,
		"Pod":            proto.EntityDTO_CONTAINER_POD,
	}

	actionTypeConverter map[string]proto.ActionItemDTO_ActionType = map[string]proto.ActionItemDTO_ActionType{
		"move": proto.ActionItemDTO_MOVE,
	}
)

func TransformActionRequest(actionAPIRequest *api.Action) (*proto.MediationServerMessage, error) {
	tSETypeString := actionAPIRequest.TargetEntityType
	glog.V(3).Infof("tType :%s", tSETypeString)
	if tSETypeString == "" {
		return nil, fmt.Errorf("Target entity type is not provided.")
	}
	targetEntityType, exist := entityTypeConverter[tSETypeString]
	if !exist {
		return nil, fmt.Errorf("Target entity type %s is not supported.", actionAPIRequest.TargetEntityType)
	}
	targetEntityDTO, err := builder.NewEntityDTOBuilder(targetEntityType,
		actionAPIRequest.TargetEntityID).Create()
	if err != nil {
		return nil, err
	}

	actionType, exist := actionTypeConverter[actionAPIRequest.ActionType]
	if !exist {
		return nil, fmt.Errorf("Action type %s is not supported.", actionType)
	}
	actionItemDTOBuilder := turbomessage.NewActionItemDTOBuilder(actionType).
		TargetSE(targetEntityDTO)
	switch actionType {
	case proto.ActionItemDTO_MOVE:
		glog.V(3).Infof("Move spec is %++v", actionAPIRequest.MoveSpec)
		nSETypeString := actionAPIRequest.MoveSpec.DestinationEntityType
		if nSETypeString == "" {
			return nil, fmt.Errorf("New service entity type is not provide for move action.")
		}
		newEntityType, exist := entityTypeConverter[nSETypeString]
		if !exist {
			return nil, fmt.Errorf("Destination entity type %s is not supported.",
				actionAPIRequest.MoveSpec.DestinationEntityType)
		}

		newEntityDTO, err := builder.NewEntityDTOBuilder(newEntityType,
			actionAPIRequest.MoveSpec.DestinationEntityID).Create()
		if err != nil {
			return nil, err
		}
		//TODO need to change go-sdk
		if newEntityType == proto.EntityDTO_VIRTUAL_MACHINE {
			virtualMachineData := &proto.EntityDTO_VirtualMachineData{
				IpAddress: []string{actionAPIRequest.MoveSpec.MoveDestinationIP},
			}
			newEntityDTO.VirtualMachineData = virtualMachineData
		}
		actionItemDTOBuilder.NewSE(newEntityDTO)
	}
	actionItemDTO, err := actionItemDTOBuilder.Build()
	if err != nil {
		return nil, err
	}
	actionExecutionDTO, err := turbomessage.NewActionExecutionDTOBuilder(actionType).
		ActionItem(actionItemDTO).Build()
	if err != nil {
		return nil, err
	}
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

	// TODO message ID is a random number in [0, 1000)
	messageID := rand.Int31n(1000)
	serverMessage := turbomessage.NewMediationServerMessageBuilder(messageID).ActionRequest(actionRequest).Build()
	return serverMessage, nil
}
