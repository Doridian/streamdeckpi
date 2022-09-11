package homeassistant

import (
	"github.com/Doridian/go-haws"
	"github.com/Doridian/streamdeckpi/agent/controller"
	"gopkg.in/yaml.v3"
)

type haEntityActionBase struct {
	haAction

	Domain string `yaml:"domain"`
	Entity string `yaml:"entity"`

	ServiceName   string                  `yaml:"service_name"`
	ServiceData   map[string]interface{}  `yaml:"service_data"`
	ServiceTarget *haws.CallServiceTarget `yaml:"service_target"`

	// TODO: Error icon and timeout
}

func (a *haEntityActionBase) ApplyConfig(config *yaml.Node, imageHelper controller.ImageHelper, ctrl controller.Controller) error {
	err := a.haAction.ApplyConfig(config, imageHelper, ctrl)
	if err != nil {
		return err
	}

	err = config.Decode(a)
	if err != nil {
		return err
	}

	if a.ServiceTarget == nil {
		a.ServiceTarget = &haws.CallServiceTarget{
			EntityID: []string{a.Entity},
		}
	}

	return nil
}

func (a *haEntityActionBase) Run(pressed bool) error {
	if !pressed {
		return nil
	}
	return a.instance.client.CallService(a.Domain, a.ServiceName, a.ServiceData, a.ServiceTarget)
}
