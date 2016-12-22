package turbomessage

import (
	"fmt"

	"github.com/vmturbo/vmturbo-go-sdk/pkg/proto"
	"crypto/rand"
)

const (
	// If simulation does not need a valid probe value, use the default one.
	DefaultProbeType string = "Default_Turbo_Probe"
)

var (
	DefaultAccountKey    string                = "Default_Turbo_Account"
	DefaultAccountValues []*proto.AccountValue = []*proto.AccountValue{
		{
			Key: &DefaultAccountKey,
			GroupScopePropertyValues: []*proto.AccountValue_PropertyValueList{},
		},
	}
)

type ActionRequestBuilder struct {
	probeType             *string
	accountValue          []*proto.AccountValue
	executionDTO          *proto.ActionExecutionDTO
	secondaryAccountValue []*proto.AccountValue

	err error
}

func NewActionRequestBuilder(probeType string, accountValue []*proto.AccountValue,
	executionDTO *proto.ActionExecutionDTO) *ActionRequestBuilder {
	var err error
	if probeType == "" {
		err = fmt.Errorf("Probe type is not provided")
	} else if executionDTO == nil {
		err = fmt.Errorf("executionDTO is not provided")
	} else if accountValue == nil {
		err = fmt.Errorf("account value is not provided.")
	}
	return &ActionRequestBuilder{
		probeType:    &probeType,
		accountValue: accountValue,
		executionDTO: executionDTO,

		err: err,
	}
}

func (arb *ActionRequestBuilder) Build() (*proto.ActionRequest, error) {
	if arb.err != nil {
		return nil, arb.err
	}
	return &proto.ActionRequest{
		ProbeType:             arb.probeType,
		AccountValue:          arb.accountValue,
		ActionExecutionDTO:    arb.executionDTO,
		SecondaryAccountValue: arb.secondaryAccountValue,
	}, nil
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
	uuid :=randUUID()
	return &ActionItemDTOBuilder{
		uuid:&uuid,
		actionType: &aType,
	}
}

func (aib *ActionItemDTOBuilder) Build() (*proto.ActionItemDTO, error) {
	if aib.err != nil {
		return nil, aib.err
	}
	return &proto.ActionItemDTO{
		ActionType:         aib.actionType,
		Uuid:               aib.uuid,
		TargetSE:           aib.targetSE,
		HostedBySE:         aib.hostedBySE,
		CurrentSE:          aib.currentSE,
		NewSE:              aib.newSE,
		CurrentComm:        aib.currentComm,
		NewComm:            aib.newComm,
		CommodityAttribute: aib.commodityAttribute,
		Providers:          aib.providers,
		EntityProfileDTO:   aib.entityProfileDTO,
		ContextData:        aib.contextData,
	}, nil
}

func (aib *ActionItemDTOBuilder) TargetSE(target *proto.EntityDTO) *ActionItemDTOBuilder {
	if aib.err != nil {
		return aib
	}
	aib.targetSE = target
	return aib
}
func (aib *ActionItemDTOBuilder) NewSE(new *proto.EntityDTO) *ActionItemDTOBuilder {
	if aib.err != nil {
		return aib
	}
	aib.newSE = new
	return aib
}


func randUUID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}