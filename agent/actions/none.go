package actions

import "github.com/Doridian/streamdeckpi/agent/interfaces"

type None struct {
	actionWithIcon
}

func (a *None) Run(pressed bool, controller interfaces.Controller) error {
	return nil
}
