package rest

import (
	//"fmt"
	"net/http"
	//"strings"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type RESTManager struct{
	handler *APIHandler

	mediationServerMessageGeneratorChan chan *proto.MediationServerMessage
}

func NewRESTManager() *RESTManager {
	return &RESTManager{
		NewAPIHandler(),
		make(chan *proto.MediationServerMessage),
	}
}

// Forward an API request to API handler. If this is a POST request, forward the generated
// MediationServerMessage to channel.
func (m *RESTManager) HandleRequest(w http.ResponseWriter, r *http.Request) {
	mediationServerMessage := m.handler.handleAPIRequest(w, r)
	if mediationServerMessage != nil && r.Method == "POST" {
		m.mediationServerMessageGeneratorChan <- mediationServerMessage
	}
}

func (m *RESTManager) ReceiveMessage() <- chan *proto.MediationServerMessage {
	return m.mediationServerMessageGeneratorChan
}