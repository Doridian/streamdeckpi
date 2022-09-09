package loader

import (
	"errors"
	"fmt"

	"github.com/Doridian/streamdeckpi/agent/actions"
	"github.com/Doridian/streamdeckpi/agent/actions/misc"
	"github.com/Doridian/streamdeckpi/agent/actions/page"
	"github.com/Doridian/streamdeckpi/agent/controller"
	"github.com/Doridian/streamdeckpi/agent/utils"
)

var actionsMap = loadActions()

func loadActions() map[string](func() actions.Action) {
	actionsMapTmp := [](func() actions.Action){
		func() actions.Action { return &misc.None{} },

		func() actions.Action { return &page.SwapPage{} },
		func() actions.Action { return &page.SwapPage{} },
		func() actions.Action { return &page.PushPage{} },
		func() actions.Action { return &page.PopPage{} },
	}

	res := make(map[string](func() actions.Action))
	for _, a := range actionsMapTmp {
		res[a().Name()] = a
	}
	return res
}

func LoadAction(name string, config *utils.YAMLRawMessage, imageLoader controller.ImageLoader, controller controller.Controller) (actions.Action, error) {
	actionCtor := actionsMap[name]
	if actionCtor == nil {
		return nil, fmt.Errorf("no action known with name: %s", name)
	}

	action := actionCtor()
	if action == nil {
		return nil, errors.New("action constructor failed")
	}

	err := config.Unmarshal(action.GetConfigRef())
	if err != nil {
		return nil, err
	}
	err = action.ApplyConfig(imageLoader, controller)
	if err != nil {
		return nil, err
	}
	return action, nil
}
