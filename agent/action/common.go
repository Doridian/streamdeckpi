package action

import (
	"github.com/Doridian/streamdeck"
	"github.com/Doridian/streamdeckpi/agent/controller"
	"gopkg.in/yaml.v3"
)

type ActionBase struct {
	ImageLoader controller.ImageLoader
	Controller  controller.Controller
}

func (a *ActionBase) ApplyConfig(config *yaml.Node, imageLoader controller.ImageLoader, ctrl controller.Controller) error {
	err := config.Decode(a)
	if err != nil {
		return err
	}

	a.ImageLoader = imageLoader
	a.Controller = ctrl
	return nil
}

type ActionWithIcon struct {
	ActionBase
	Icon string `yaml:"icon"`
}

func (a *ActionWithIcon) ApplyConfig(config *yaml.Node, imageLoader controller.ImageLoader, ctrl controller.Controller) error {
	err := a.ActionBase.ApplyConfig(config, imageLoader, ctrl)
	if err != nil {
		return err
	}
	return config.Decode(a)
}

func (a *ActionWithIcon) Render(force bool) (*streamdeck.ImageData, error) {
	if force {
		return a.ImageLoader.Load(a.Icon)
	}
	return nil, nil
}
