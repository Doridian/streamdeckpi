package homeassistant

import (
	"github.com/Doridian/go-haws"
	"github.com/Doridian/streamdeck"
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
	DefaultIcon string `yaml:"default_icon"`

	currentIcon      string
	lastRenderedIcon string
}

func (a *haEntityActionBase) ApplyConfig(config *yaml.Node, imageLoader controller.ImageLoader, ctrl controller.Controller) error {
	err := a.haAction.ApplyConfig(config, imageLoader, ctrl)
	if err != nil {
		return err
	}

	err = config.Decode(a)
	if err != nil {
		return err
	}

	a.currentIcon = a.DefaultIcon
	if a.ServiceTarget == nil {
		a.ServiceTarget = &haws.CallServiceTarget{
			EntityID: []string{a.Entity},
		}
	}

	return nil
}

func (a *haEntityActionBase) Render(force bool) (*streamdeck.ImageData, error) {
	if a.currentIcon == a.lastRenderedIcon && !force {
		return nil, nil
	}

	a.lastRenderedIcon = a.currentIcon
	return a.ImageLoader.Load(a.lastRenderedIcon)
}

func (a *haEntityActionBase) Run(pressed bool) error {
	if !pressed {
		return nil
	}
	return a.instance.client.CallService(a.Domain, a.ServiceName, a.ServiceData, a.ServiceTarget)
}
