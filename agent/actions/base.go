package actions

import (
	"errors"
	"image"

	"github.com/Doridian/streamdeckpi/agent/interfaces"
)

type Action interface {
	SetConfig(config map[string]interface{}) error

	Run(pressed bool, controller interfaces.Controller) error

	// In Render, you can return a nil image to indicate the image hasn't changed since the last call
	// This will indicate to the renderer to not change the image
	//
	// If force is true, you must always return an image if the action has one available
	// Otherwsie, a blank image will be set
	Render(force bool) (image.Image, error)
}

func LoadAction(name string, config map[string]interface{}) (Action, error) {
	return nil, errors.New("implement me")
}
