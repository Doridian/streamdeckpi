package page

import (
	"github.com/Doridian/streamdeckpi/agent/action"
	"github.com/Doridian/streamdeckpi/agent/controller"
	"gopkg.in/yaml.v3"
)

type PushPage struct {
	action.ActionWithIcon
	Target string `yaml:"target"`
}

func (a *PushPage) ApplyConfig(config *yaml.Node, imageLoader controller.ImageLoader, ctrl controller.Controller) error {
	err := a.ActionWithIcon.ApplyConfig(config, imageLoader, ctrl)
	if err != nil {
		return err
	}
	return config.Decode(a)
}

func (a *PushPage) Run(pressed bool) error {
	if !pressed {
		return nil
	}
	return a.Controller.PushPage(a.Target)
}

func (a *PushPage) Name() string {
	return "push_page"
}
