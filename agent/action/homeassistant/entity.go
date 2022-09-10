package homeassistant

import (
	"github.com/Doridian/go-haws"
	"github.com/Doridian/streamdeckpi/agent/action"
	"github.com/Doridian/streamdeckpi/agent/controller"
	"gopkg.in/yaml.v3"
)

type haConditionalIcons struct {
	Condition haCondition `yaml:"condition"`
	File      string      `yaml:"file"`
}

type haEntityAction struct {
	haEntityActionBase

	Icons []haConditionalIcons `yaml:"icons"`
}

func (a *haEntityAction) New() action.Action {
	return &haEntityAction{}
}

func (a *haEntityAction) OnState(entityID string, state haws.State) error {
	foundIcon := ""

	for _, icon := range a.Icons {
		match, err := icon.Condition.Evaluate(state.State)
		if err != nil {
			return err
		}

		if match {
			foundIcon = icon.File
			break
		}
	}

	if foundIcon == "" {
		foundIcon = a.DefaultIcon
	}

	a.currentIcon = foundIcon
	return nil
}

func (a *haEntityAction) ApplyConfig(config *yaml.Node, imageLoader controller.ImageLoader, ctrl controller.Controller) error {
	err := a.haEntityActionBase.ApplyConfig(config, imageLoader, ctrl)
	if err != nil {
		return err
	}

	err = config.Decode(a)
	if err != nil {
		return err
	}

	if a.Icons == nil {
		a.Icons = make([]haConditionalIcons, 0)
	}

	a.instance.RegisterStateReceiver(a, a.Entity)

	return nil
}

func (a *haEntityAction) Name() string {
	return "homeassistant_entity"
}
