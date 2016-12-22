package turbomessage

import (
	"fmt"

	"github.com/vmturbo/vmturbo-go-sdk/pkg/proto"
)

type ActionRequestBuilder struct {
	probeType             *string
	accountValue          []*proto.AccountValue
	executionDTO          *proto.ActionExecutionDTO
	secondaryAccountValue []*proto.AccountValue
}

func NewActionRequestBuilder(executionDTO *proto.ActionExecutionDTO) *ActionRequestBuilder {
	return &ActionRequestBuilder{
		executionDTO: executionDTO,
	}
}

func (arb *ActionRequestBuilder) Build() *proto.ActionRequest {
	return &proto.ActionRequest{
		ProbeType:             arb.probeType,
		AccountValue:          arb.accountValue,
		ActionExecutionDTO:    arb.executionDTO,
		SecondaryAccountValue: arb.secondaryAccountValue,
	}
}

type ActionExecutionDTOBuilder struct {
	actionType *proto.ActionItemDTO_ActionType
	actionItem []*proto.ActionItemDTO
	progress   *int64

	err error
}

func NewActionExecutionDTOBuilder(aType proto.ActionItemDTO_ActionType) *ActionExecutionDTOBuilder {
	return &ActionExecutionDTOBuilder{
		actionType: &aType,
	}
}

func (aeb *ActionExecutionDTOBuilder) Build() (*proto.ActionExecutionDTO, error) {
	if aeb.err != nil {
		return nil, aeb.err
	}
	return &proto.ActionExecutionDTO{
		ActionType: aeb.actionType,
		ActionItem: aeb.actionItem,
		Progress:   aeb.progress,
	}, nil
}

func (aeb *ActionExecutionDTOBuilder) ActionItem(actionItem *proto.ActionItemDTO) *ActionExecutionDTOBuilder {
	if aeb.err != nil {
		return aeb
	}
	if actionItem == nil {
		aeb.err = fmt.Errorf("ActionItem passed in is nil")
		return aeb
	}
	if aeb.actionItem == nil {
		aeb.actionItem = []*proto.ActionItemDTO{}
	}
	aeb.actionItem = append(aeb.actionItem, actionItem)
	return aeb
}

type ActionItemDTOBuilder struct {
	actionType         *proto.ActionItemDTO_ActionType
	uuid               *string
	targetSE           *proto.EntityDTO
	hostedBySE         *proto.EntityDTO
	currentSE          *proto.EntityDTO
	newSE              *proto.EntityDTO
	currentComm        *proto.CommodityDTO
	newComm            *proto.CommodityDTO
	commodityAttribute *proto.ActionItemDTO_CommodityAttribute
	providers          []*proto.ActionItemDTO_ProviderInfo
	entityProfileDTO   *proto.EntityProfileDTO
	contextData        []*proto.ContextData

	err error
}

func NewActionItemDTOBuilder(aType proto.ActionItemDTO_ActionType) *ActionItemDTOBuilder {
	return &ActionItemDTOBuilder{
		actionType: &aType,
	}
}

func (aib *ActionItemDTOBuilder) Build() (*proto.ActionItemDTO, error) {
	if aib.err != nil {
		return nil, aib.err
	}
	return &proto.ActionItemDTO {
		ActionType:aib.actionType,
		Uuid: aib.uuid,
		TargetSE:aib.targetSE,
		HostedBySE:aib.hostedBySE,
		CurrentSE:aib.currentSE          ,
		NewSE:aib.newSE,
		CurrentComm:aib.currentComm,
		NewComm:aib.newComm,
		CommodityAttribute: aib.commodityAttribute,
		Providers:aib.providers,
		EntityProfileDTO:aib.entityProfileDTO,
		ContextData:aib.contextData,
	}, nil
}

func (aib *ActionItemDTOBuilder) TargetSE(target *proto.EntityDTO) *ActionItemDTOBuilder {
	if (aib.err != nil) {
		return aib
	}
	aib.targetSE = target
	return aib
}
func (aib *ActionItemDTOBuilder) NewSE(new *proto.EntityDTO) *ActionItemDTOBuilder {
	if (aib.err != nil) {
		return aib
	}
	aib.newSE = new
	return aib
}
