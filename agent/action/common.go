package action

import (
	"github.com/Doridian/go-streamdeck"
	"github.com/Doridian/streamdeckpi/agent/controller"
	"gopkg.in/yaml.v3"
)

type ActionBase struct {
	ImageHelper controller.ImageHelper
	Controller  controller.Controller
}

func (a *ActionBase) ApplyConfig(config *yaml.Node, imageHelper controller.ImageHelper, ctrl controller.Controller) error {
	err := config.Decode(a)
	if err != nil {
		return err
	}

	a.ImageHelper = imageHelper
	a.Controller = ctrl
	return nil
}

type ActionWithIcon struct {
	ActionBase
	Icon string `yaml:"icon"`
}

func (a *ActionWithIcon) ApplyConfig(config *yaml.Node, imageHelper controller.ImageHelper, ctrl controller.Controller) error {
	err := a.ActionBase.ApplyConfig(config, imageHelper, ctrl)
	if err != nil {
		return err
	}
	return config.Decode(a)
}

func (a *ActionWithIcon) Render(force bool) (*streamdeck.ImageData, error) {
	if force {
		return a.ImageHelper.Load(a.Icon)
	}
	return nil, nil
}
