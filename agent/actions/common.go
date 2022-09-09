package actions

import (
	"github.com/Doridian/streamdeck"
	"github.com/Doridian/streamdeckpi/agent/controller"
)

type ActionBase struct {
	ImageLoader controller.ImageLoader
	Controller  controller.Controller
}

func (a *ActionBase) ApplyConfig(imageLoader controller.ImageLoader, controller controller.Controller) error {
	a.ImageLoader = imageLoader
	a.Controller = controller
	return nil
}

type ActionWithIcon struct {
	ActionBase
	Icon string `yaml:"icon"`
}

func (a *ActionWithIcon) Render(force bool) (*streamdeck.ImageData, error) {
	if force {
		return a.ImageLoader.Load(a.Icon)
	}
	return nil, nil
}
