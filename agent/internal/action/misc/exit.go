package misc

import (
	"os"

	"github.com/Doridian/streamdeckpi/agent/internal/action"
)

type exit struct {
	action.ActionWithIcon
}

func (a *exit) New() action.Action {
	return &exit{}
}

func (a *exit) Run(pressed bool) error {
	if !pressed {
		return nil
	}
	os.Exit(0)
	return nil
}

func (a *exit) Name() string {
	return "exit"
}
