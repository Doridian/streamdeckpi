package action

import (
	"github.com/Doridian/go-streamdeck"
	"github.com/Doridian/streamdeckpi/agent/internal/controller"
	"gopkg.in/yaml.v3"
)

type Action interface {
	New() Action

	// When ApplyConfig is called, apply config
	// Your struct is expected to declare exported fields with YAML annotation
	// for config values
	ApplyConfig(config *yaml.Node, imageHelper controller.ImageHelper, ctrl controller.Controller) error

	Run(pressed bool) error
	Name() string

	// In Render, you can return a nil image to indicate the image hasn't changed since the last call
	// This will indicate to the renderer to not change the image
	//
	// If force is true, you must always return an image if the action has one available
	// Otherwsie, a blank image will be set
	Render(force bool) (*streamdeck.ImageData, error)
}
