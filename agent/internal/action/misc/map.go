package misc

import (
	"github.com/Doridian/go-streamdeck"
	"github.com/Doridian/streamdeckpi/agent/internal/action"
	"github.com/Doridian/streamdeckpi/agent/internal/action/loader"
	"github.com/Doridian/streamdeckpi/agent/internal/controller"
	"gopkg.in/yaml.v3"
)

type mapActionConfig struct {
	Name       string    `yaml:"name"`
	Parameters yaml.Node `yaml:"parameters"`
}

type reMap struct {
	action.ActionBase

	RunActionConfig    *mapActionConfig `yaml:"run"`
	RenderActionConfig *mapActionConfig `yaml:"render"`

	runAction    action.Action
	renderAction action.Action
}

func (a *reMap) New() action.Action {
	return &reMap{}
}

func (a *reMap) ApplyConfig(config *yaml.Node, imageHelper controller.ImageHelper, ctrl controller.Controller) error {
	err := a.ActionBase.ApplyConfig(config, imageHelper, ctrl)
	if err != nil {
		return err
	}

	err = config.Decode(a)
	if err != nil {
		return err
	}

	a.runAction, err = loader.LoadAction(a.RunActionConfig.Name, &a.RunActionConfig.Parameters, imageHelper, ctrl)
	if err != nil {
		return err
	}

	a.renderAction, err = loader.LoadAction(a.RenderActionConfig.Name, &a.RenderActionConfig.Parameters, imageHelper, ctrl)
	if err != nil {
		return err
	}

	return nil
}

func (a *reMap) Run(pressed bool) error {
	return a.runAction.Run(pressed)
}

func (a *reMap) Render(force bool) (*streamdeck.ImageData, error) {
	return a.renderAction.Render(force)
}

func (a *reMap) Name() string {
	return "remap"
}
