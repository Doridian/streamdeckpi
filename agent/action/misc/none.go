package misc

import "github.com/Doridian/streamdeckpi/agent/action"

type none struct {
	action.ActionWithIcon
}

func (a *none) New() action.Action {
	return &none{}
}

func (a *none) Run(pressed bool) error {
	return nil
}

func (a *none) Name() string {
	return "none"
}
