package homeassistant

import (
	"github.com/Doridian/go-haws"
	"github.com/Doridian/go-streamdeck"
	"github.com/Doridian/streamdeckpi/agent/action"
	"github.com/Doridian/streamdeckpi/agent/controller"
	"gopkg.in/yaml.v3"
)

type haConditionOverride struct {
	Condition haCondition `yaml:"condition"`
	Icon      string      `yaml:"icon"`

	Domain        string                  `yaml:"domain"`
	ServiceName   string                  `yaml:"service_name"`
	ServiceData   map[string]interface{}  `yaml:"service_data"`
	ServiceTarget *haws.CallServiceTarget `yaml:"service_target"`
}

type haEntityAction struct {
	haEntityActionBase
	Icon string `yaml:"icon"`

	Conditions []*haConditionOverride `yaml:"conditions"`

	currentDomain        string
	currentServiceName   string
	currentServiceData   map[string]interface{}
	currentServiceTarget *haws.CallServiceTarget

	currentIcon      string
	lastRenderedIcon string
}

func (a *haEntityAction) New() action.Action {
	return &haEntityAction{}
}

func (a *haEntityAction) OnState(entityID string, state haws.State) error {
	var currentMatch *haConditionOverride
	for _, cond := range a.Conditions {
		match, err := cond.Condition.Evaluate(&state)
		if err != nil {
			return err
		}

		if match {
			currentMatch = cond
			break
		}
	}

	foundIcon := ""
	foundDomain := ""
	foundServiceName := ""
	var foundServiceData map[string]interface{}
	var foundServiceTarget *haws.CallServiceTarget

	if currentMatch != nil {
		foundIcon = currentMatch.Icon
		foundDomain = currentMatch.Domain
		foundServiceName = currentMatch.ServiceName
		foundServiceData = currentMatch.ServiceData
		foundServiceTarget = currentMatch.ServiceTarget
	}

	if foundIcon == "" {
		foundIcon = a.Icon
	}

	if foundDomain == "" {
		foundDomain = a.Domain
	}

	if foundServiceName == "" {
		foundServiceName = a.ServiceName
	}

	if foundServiceData == nil {
		foundServiceData = a.ServiceData
	}

	if foundServiceTarget == nil {
		foundServiceTarget = a.ServiceTarget
	}

	a.currentIcon = foundIcon
	a.currentDomain = foundDomain
	a.currentServiceName = foundServiceName
	a.currentServiceData = foundServiceData
	a.currentServiceTarget = foundServiceTarget

	return nil
}

func (a *haEntityAction) ApplyConfig(config *yaml.Node, imageHelper controller.ImageHelper, ctrl controller.Controller) error {
	err := a.haEntityActionBase.ApplyConfig(config, imageHelper, ctrl)
	if err != nil {
		return err
	}

	err = config.Decode(a)
	if err != nil {
		return err
	}

	if a.Conditions == nil {
		a.Conditions = make([]*haConditionOverride, 0)
	}

	a.instance.RegisterStateReceiver(a, a.Entity)

	return nil
}

func (a *haEntityAction) Name() string {
	return "homeassistant_entity"
}

func (a *haEntityAction) Run(pressed bool) error {
	if !pressed {
		return nil
	}
	return a.instance.client.CallService(a.currentDomain, a.currentServiceName, a.currentServiceData, a.currentServiceTarget)
}

func (a *haEntityAction) Render(force bool) (*streamdeck.ImageData, error) {
	toRender := a.currentIcon
	if toRender == a.lastRenderedIcon && !force {
		return nil, nil
	}

	a.lastRenderedIcon = toRender
	return a.ImageHelper.Load(toRender)
}
