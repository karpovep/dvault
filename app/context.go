package app

import "dvault/constants"

type (
	IAppContext interface {
		Set(name constants.AppContextKey, resource interface{})
		Get(name constants.AppContextKey) interface{}
	}

	AppContext struct {
		resources map[constants.AppContextKey]interface{}
	}
)

func NewApplicationContext() IAppContext {
	return &AppContext{
		resources: make(map[constants.AppContextKey]interface{}),
	}
}

func (ac *AppContext) Set(name constants.AppContextKey, resource interface{}) {
	ac.resources[name] = resource
}

func (ac *AppContext) Get(name constants.AppContextKey) interface{} {
	return ac.resources[name]
}
