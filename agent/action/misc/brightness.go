package misc

import (
	"github.com/Doridian/streamdeckpi/agent/action"
	"github.com/Doridian/streamdeckpi/agent/controller"
	"gopkg.in/yaml.v3"
)

type Brightness struct {
	action.ActionWithIcon

	Brightness int `yaml:"brightness"`
}

func (a *Brightness) ApplyConfig(config *yaml.Node, imageLoader controller.ImageLoader, ctrl controller.Controller) error {
	err := a.ActionWithIcon.ApplyConfig(config, imageLoader, ctrl)
	if err != nil {
		return err
	}
	return config.Decode(a)
}

func (a *Brightness) Run(pressed bool) error {
	if !pressed {
		return nil
	}

	return a.Controller.SetBrightness(a.Brightness)
}

func (a *Brightness) Name() string {
	return "brightness"
}
