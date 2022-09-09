package misc

import (
	"os"

	"github.com/Doridian/streamdeckpi/agent/actions"
)

type Exit struct {
	actions.ActionWithIcon
}

func (a *Exit) Run(pressed bool) error {
	if !pressed {
		return nil
	}
	os.Exit(0)
	return nil
}

func (a *Exit) Name() string {
	return "exit"
}
