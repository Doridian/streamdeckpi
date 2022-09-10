package page

import (
	"github.com/Doridian/streamdeckpi/agent/action"
	"github.com/Doridian/streamdeckpi/agent/controller"
	"gopkg.in/yaml.v3"
)

type swapPage struct {
	action.ActionWithIcon
	Target string `yaml:"target"`
}

func (a *swapPage) New() action.Action {
	return &swapPage{}
}

func (a *swapPage) ApplyConfig(config *yaml.Node, imageLoader controller.ImageLoader, controller controller.Controller) error {
	err := a.ActionWithIcon.ApplyConfig(config, imageLoader, controller)
	if err != nil {
		return err
	}
	return config.Decode(a)
}

func (a *swapPage) Run(pressed bool) error {
	if !pressed {
		return nil
	}
	return a.Controller.SwapPage(a.Target)
}

func (a *swapPage) Name() string {
	return "swap_page"
}
