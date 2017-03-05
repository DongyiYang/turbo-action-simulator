package rest

import (
	"net/http"

	"github.com/turbonomic/turbo-simulator/pkg/rest/api"
)

type RESTManager struct {
	handler *APIHandler

	apiObjectGeneratorChan chan api.APIObject
}

func NewRESTManager() *RESTManager {
	apiObjectGeneratorChan := make(chan api.APIObject)
	return &RESTManager{
		NewAPIHandler(apiObjectGeneratorChan),
		apiObjectGeneratorChan,
	}
}

// Forward an API request to API handler.
// If this is a POST request, forward the generated MediationServerMessage to channel.
func (m *RESTManager) HandleRequest(w http.ResponseWriter, r *http.Request) {
	m.handler.handleAPIRequest(w, r)
}

func (m *RESTManager) ReceiveMessage() <-chan api.APIObject {
	return m.apiObjectGeneratorChan
}
