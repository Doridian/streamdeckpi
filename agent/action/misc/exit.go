package misc

import (
	"os"

	"github.com/Doridian/streamdeckpi/agent/action"
)

type Exit struct {
	action.ActionWithIcon
}

func (a *Exit) New() action.Action {
	return &Exit{}
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
