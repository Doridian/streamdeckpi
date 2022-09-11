package homeassistant

import (
	"github.com/Doridian/streamdeckpi/agent/action"
	"github.com/Doridian/streamdeckpi/agent/controller"
	"gopkg.in/yaml.v3"
)

type haAction struct {
	action.ActionBase

	Instance string `yaml:"instance"`
	instance *haInstance
}

func (a *haAction) ApplyConfig(config *yaml.Node, imageHelper controller.ImageHelper, ctrl controller.Controller) error {
	err := a.ActionBase.ApplyConfig(config, imageHelper, ctrl)
	if err != nil {
		return err
	}

	err = config.Decode(a)
	if err != nil {
		return err
	}

	a.instance, err = GetHomeAssistant(ctrl, a.Instance)
	if err != nil {
		return err
	}

	return nil
}
