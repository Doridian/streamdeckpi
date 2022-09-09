package misc

import "github.com/Doridian/streamdeckpi/agent/action"

type Reset struct {
	action.ActionWithIcon
}

func (a *Reset) Run(pressed bool) error {
	return a.Controller.Reset()
}

func (a *Reset) Name() string {
	return "reset"
}
