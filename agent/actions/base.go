package actions

import (
	"errors"
	"fmt"

	"github.com/Doridian/streamdeck"
	"github.com/Doridian/streamdeckpi/agent/interfaces"
	"github.com/Doridian/streamdeckpi/agent/utils"
)

type Action interface {
	// Return a reference to your config object in GetConfigRef
	// When ApplyConfig is called, read whatever is in it (it will be written to by the config loader)
	GetConfigRef() interface{}
	ApplyConfig(imageLoader interfaces.ImageLoader, controller interfaces.Controller) error

	Run(pressed bool) error
	Name() string

	// In Render, you can return a nil image to indicate the image hasn't changed since the last call
	// This will indicate to the renderer to not change the image
	//
	// If force is true, you must always return an image if the action has one available
	// Otherwsie, a blank image will be set
	Render(force bool) (*streamdeck.ImageData, error)
}

var actionsMap = loadActions()

func loadActions() map[string](func() Action) {
	actions := [](func() Action){
		func() Action { return &None{} },

		func() Action { return &SwapPage{} },
		func() Action { return &SwapPage{} },
		func() Action { return &PushPage{} },
		func() Action { return &PopPage{} },
	}

	res := make(map[string](func() Action))
	for _, a := range actions {
		res[a().Name()] = a
	}
	return res
}

func LoadAction(name string, config *utils.YAMLRawMessage, imageLoader interfaces.ImageLoader, controller interfaces.Controller) (Action, error) {
	actionCtor := actionsMap[name]
	if actionCtor == nil {
		return nil, fmt.Errorf("no action known with name: %s", name)
	}

	action := actionCtor()
	if action == nil {
		return nil, errors.New("action constructor failed")
	}

	err := config.Unmarshal(action.GetConfigRef())
	if err != nil {
		return nil, err
	}
	err = action.ApplyConfig(imageLoader, controller)
	if err != nil {
		return nil, err
	}
	return action, nil
}
