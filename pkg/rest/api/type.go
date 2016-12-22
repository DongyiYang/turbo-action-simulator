package api

type Action struct {
	ActionType       string    `json:"type,omitempty"`
	TargetEntityType string    `json:"entityType,omitempty"`
	TargetEntityID   string    `json:"entityID,omitempty"`
	MoveSpec         *MoveSpec `json:"moveSpec,omitempty"`
}

type MoveSpec struct {
	DestinationEntityType string `destinationEntityType,omitempty`
	DestinationEntityID   string `destinationEntityID,omitempty`
	MoveDestinationIP     string `json:"destinationIP,omitempty"`
}
