package loader

import (
	"errors"
	"fmt"

	"github.com/Doridian/streamdeckpi/agent/action"
	"github.com/Doridian/streamdeckpi/agent/controller"
	"gopkg.in/yaml.v3"
)

var actionsMap = make(map[string]action.Action)

func RegisterAction(impl action.Action) {
	actionsMap[impl.Name()] = impl
}

func LoadAction(name string, config *yaml.Node, imageHelper controller.ImageHelper, ctrl controller.Controller) (action.Action, error) {
	actionCtor := actionsMap[name]
	if actionCtor == nil {
		return nil, fmt.Errorf("no action known with name: %s", name)
	}

	actionObj := actionCtor.New()
	if actionObj == nil {
		return nil, errors.New("action constructor failed")
	}

	err := actionObj.ApplyConfig(config, imageHelper, ctrl)
	if err != nil {
		return nil, err
	}
	return actionObj, nil
}
