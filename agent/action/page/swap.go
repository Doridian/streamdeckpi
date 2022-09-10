package page

import (
	"github.com/Doridian/streamdeckpi/agent/action"
	"github.com/Doridian/streamdeckpi/agent/controller"
	"gopkg.in/yaml.v3"
)

type SwapPage struct {
	action.ActionWithIcon
	Target string `yaml:"target"`
}

func (a *SwapPage) New() action.Action {
	return &SwapPage{}
}

func (a *SwapPage) ApplyConfig(config *yaml.Node, imageLoader controller.ImageLoader, controller controller.Controller) error {
	err := a.ActionWithIcon.ApplyConfig(config, imageLoader, controller)
	if err != nil {
		return err
	}
	return config.Decode(a)
}

func (a *SwapPage) Run(pressed bool) error {
	if !pressed {
		return nil
	}
	return a.Controller.SwapPage(a.Target)
}

func (a *SwapPage) Name() string {
	return "swap_page"
}
