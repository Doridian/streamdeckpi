package actions

import (
	"fmt"

	"github.com/Doridian/streamdeck"
	"github.com/Doridian/streamdeckpi/agent/interfaces"
	"github.com/Doridian/streamdeckpi/agent/utils"
)

type Action interface {
	// Return a reference to your config object in GetConfigRef
	// When ApplyConfig is called, read whatever is in it (it will be written to by the config loader)
	GetConfigRef() interface{}
	ApplyConfig(imageLoader interfaces.ImageLoader) error

	Run(pressed bool, controller interfaces.Controller) error

	// In Render, you can return a nil image to indicate the image hasn't changed since the last call
	// This will indicate to the renderer to not change the image
	//
	// If force is true, you must always return an image if the action has one available
	// Otherwsie, a blank image will be set
	Render(force bool) (*streamdeck.ImageData, error)
}

func LoadAction(name string, config *utils.YAMLRawMessage, imageLoader interfaces.ImageLoader) (Action, error) {
	var action Action
	// TODO: Find action struct and create here
	if action == nil {
		return nil, fmt.Errorf("no action known with name: %s", name)
	}

	err := config.Unmarshal(action.GetConfigRef())
	if err != nil {
		return nil, err
	}
	err = action.ApplyConfig(imageLoader)
	if err != nil {
		return nil, err
	}
	return action, nil
}
