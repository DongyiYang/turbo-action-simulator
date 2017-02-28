package rest

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"net/http"
)

type RESTManager struct {
	handler *APIHandler

	mediationServerMessageGeneratorChan chan *proto.MediationServerMessage
}

func NewRESTManager() *RESTManager {
	mediationServerMessageGeneratorChan := make(chan *proto.MediationServerMessage)
	return &RESTManager{
		NewAPIHandler(mediationServerMessageGeneratorChan),
		mediationServerMessageGeneratorChan,
	}
}

// Forward an API request to API handler.
// If this is a POST request, forward the generated MediationServerMessage to channel.
func (m *RESTManager) HandleRequest(w http.ResponseWriter, r *http.Request) {
	m.handler.handleAPIRequest(w, r)
}

func (m *RESTManager) ReceiveMessage() <-chan *proto.MediationServerMessage {
	return m.mediationServerMessageGeneratorChan
}
