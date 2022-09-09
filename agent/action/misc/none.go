package misc

import "github.com/Doridian/streamdeckpi/agent/action"

type None struct {
	action.ActionWithIcon
}

func (a *None) Run(pressed bool) error {
	return nil
}

func (a *None) Name() string {
	return "none"
}
