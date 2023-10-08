package page

import (
	"github.com/Doridian/streamdeckpi/agent/internal/action"
	"github.com/Doridian/streamdeckpi/agent/internal/controller"
	"gopkg.in/yaml.v3"
)

type pushPage struct {
	action.ActionWithIcon
	Target string `yaml:"target"`
}

func (a *pushPage) New() action.Action {
	return &pushPage{}
}

func (a *pushPage) ApplyConfig(config *yaml.Node, imageHelper controller.ImageHelper, ctrl controller.Controller) error {
	err := a.ActionWithIcon.ApplyConfig(config, imageHelper, ctrl)
	if err != nil {
		return err
	}
	err = config.Decode(a)
	if err != nil {
		return err
	}
	return ctrl.PreloadPage(a.Target)
}

func (a *pushPage) Run(pressed bool) error {
	if !pressed {
		return nil
	}
	return a.Controller.PushPage(a.Target)
}

func (a *pushPage) Name() string {
	return "push_page"
}
