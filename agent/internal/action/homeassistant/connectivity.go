package homeassistant

import (
	"github.com/Doridian/go-streamdeck"
	"github.com/Doridian/streamdeckpi/agent/internal/action"
	"github.com/Doridian/streamdeckpi/agent/internal/controller"
	"gopkg.in/yaml.v3"
)

type haConnectivityAction struct {
	haAction

	DefaultIcon string                     `yaml:"default_icon"`
	Icons       map[ConnectionState]string `yaml:"icons"`

	lastState ConnectionState
}

func (a *haConnectivityAction) New() action.Action {
	return &haConnectivityAction{}
}

func (a *haConnectivityAction) Name() string {
	return "homeassistant_connectivity"
}

func (a *haConnectivityAction) Run(pressed bool) error {
	return nil
}

func (a *haConnectivityAction) ApplyConfig(config *yaml.Node, imageHelper controller.ImageHelper, ctrl controller.Controller) error {
	err := a.haAction.ApplyConfig(config, imageHelper, ctrl)
	if err != nil {
		return err
	}

	return config.Decode(a)
}

func (a *haConnectivityAction) Render(force bool) (*streamdeck.ImageData, error) {
	curState := a.instance.ConnectionState()
	if curState == a.lastState && !force {
		return nil, nil
	}

	toRender := a.Icons[curState]
	if toRender == "" {
		toRender = a.DefaultIcon
		if toRender == "" {
			return nil, nil
		}
	}

	a.lastState = curState
	return a.ImageHelper.Load(toRender)
}
