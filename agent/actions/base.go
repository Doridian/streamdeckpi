package actions

import (
	"github.com/Doridian/streamdeck"
	"github.com/Doridian/streamdeckpi/agent/controller"
)

type Action interface {
	// Return a reference to your config object in GetConfigRef
	// When ApplyConfig is called, read whatever is in it (it will be written to by the config loader)
	GetConfigRef() interface{}
	ApplyConfig(imageLoader controller.ImageLoader, controller controller.Controller) error

	Run(pressed bool) error
	Name() string

	// In Render, you can return a nil image to indicate the image hasn't changed since the last call
	// This will indicate to the renderer to not change the image
	//
	// If force is true, you must always return an image if the action has one available
	// Otherwsie, a blank image will be set
	Render(force bool) (*streamdeck.ImageData, error)
}
