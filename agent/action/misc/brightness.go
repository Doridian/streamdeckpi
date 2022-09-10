package misc

import (
	"github.com/Doridian/streamdeckpi/agent/action"
	"github.com/Doridian/streamdeckpi/agent/controller"
	"gopkg.in/yaml.v3"
)

type brightness struct {
	action.ActionWithIcon

	Brightness int `yaml:"brightness"`
}

func (a *brightness) New() action.Action {
	return &brightness{}
}

func (a *brightness) ApplyConfig(config *yaml.Node, imageLoader controller.ImageLoader, ctrl controller.Controller) error {
	err := a.ActionWithIcon.ApplyConfig(config, imageLoader, ctrl)
	if err != nil {
		return err
	}
	return config.Decode(a)
}

func (a *brightness) Run(pressed bool) error {
	if !pressed {
		return nil
	}

	return a.Controller.SetBrightness(a.Brightness)
}

func (a *brightness) Name() string {
	return "brightness"
}
