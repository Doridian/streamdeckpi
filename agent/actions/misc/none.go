package misc

import "github.com/Doridian/streamdeckpi/agent/actions"

type None struct {
	actions.ActionWithIcon
}

func (a *None) Run(pressed bool) error {
	return nil
}

func (a *None) Name() string {
	return "none"
}
