package loader

import (
	"errors"
	"fmt"

	"github.com/Doridian/streamdeckpi/agent/action"
	"github.com/Doridian/streamdeckpi/agent/action/homeassistant"
	"github.com/Doridian/streamdeckpi/agent/action/misc"
	"github.com/Doridian/streamdeckpi/agent/action/page"
	"github.com/Doridian/streamdeckpi/agent/controller"
	"gopkg.in/yaml.v3"
)

var actionsMap = loadActions()

func loadActions() map[string](func() action.Action) {
	actionsList := [](func() action.Action){
		func() action.Action { return &misc.None{} },
		func() action.Action { return &misc.Exit{} },
		func() action.Action { return &misc.Command{} },
		func() action.Action { return &misc.Brightness{} },

		func() action.Action { return &page.SwapPage{} },
		func() action.Action { return &page.SwapPage{} },
		func() action.Action { return &page.PushPage{} },
		func() action.Action { return &page.PopPage{} },

		func() action.Action { return &homeassistant.HAEntityAction{} },
	}

	res := make(map[string](func() action.Action))
	for _, a := range actionsList {
		res[a().Name()] = a
	}
	return res
}

func LoadAction(name string, config *yaml.Node, imageLoader controller.ImageLoader, ctrl controller.Controller) (action.Action, error) {
	actionCtor := actionsMap[name]
	if actionCtor == nil {
		return nil, fmt.Errorf("no action known with name: %s", name)
	}

	actionObj := actionCtor()
	if actionObj == nil {
		return nil, errors.New("action constructor failed")
	}

	err := actionObj.ApplyConfig(config, imageLoader, ctrl)
	if err != nil {
		return nil, err
	}
	return actionObj, nil
}
