package misc

import "github.com/Doridian/streamdeckpi/agent/actions"

type Reset struct {
	actions.ActionWithIcon
}

func (a *Reset) Run(pressed bool) error {
	return a.Controller.Reset()
}

func (a *Reset) Name() string {
	return "reset"
}
