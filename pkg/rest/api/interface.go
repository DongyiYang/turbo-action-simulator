package api

type APIObject interface {
	IsAPIObject()
}

// register
func (Action) IsAPIObject()           {}
func (Discovery) IsAPIObject() {}
