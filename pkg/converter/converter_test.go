package converter

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbo-go-sdk/util/rand"
	"github.com/turbonomic/turbo-simulator/pkg/rest/api"
	"reflect"
	"testing"
)

func TestGetMovingTarget(t *testing.T) {
	table := []struct {
		DestinationEntityType string
		DestinationEntityID   string
		MoveDestinationIP     string

		expectErr bool
	}{
		{
			// Destination type is not provided.
			DestinationEntityType: "",

			expectErr: true,
		},
		{
			// Unknown new entity type.
			DestinationEntityType: "#" + rand.String(5),

			expectErr: true,
		},
		{
			DestinationEntityType: "VirtualMachine",
			DestinationEntityID:   rand.String(10),
			MoveDestinationIP:     "1.1.1.1",

			expectErr: false,
		},
	}

	for _, item := range table {
		moveSpec := api.MoveSpec{
			DestinationEntityType: item.DestinationEntityType,
			DestinationEntityID:   item.DestinationEntityID,
			MoveDestinationIP:     item.MoveDestinationIP,
		}
		newEntityDTO, err := getMovingTarget(moveSpec)
		if item.expectErr {
			if err == nil {
				t.Error("Expects error, got no error")
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error %s", err)
			}
			newEntityType := entityTypeConverter[item.DestinationEntityType]
			expectedEntityDTO := &proto.EntityDTO{
				Id:         &item.DestinationEntityID,
				EntityType: &newEntityType,
			}
			switch newEntityType {
			case proto.EntityDTO_VIRTUAL_MACHINE:
				virtualMachineData := &proto.EntityDTO_VirtualMachineData{
					IpAddress: []string{item.MoveDestinationIP},
				}
				expectedEntityDTO.EntityData = &proto.EntityDTO_VirtualMachineData_{
					VirtualMachineData: virtualMachineData,
				}
			}
			if !reflect.DeepEqual(expectedEntityDTO, newEntityDTO) {
				t.Errorf("Expected %v, got %v", expectedEntityDTO, newEntityDTO)
			}
		}
	}
}
