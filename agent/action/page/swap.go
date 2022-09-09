package page

import "github.com/Doridian/streamdeckpi/agent/action"

type SwapPage struct {
	action.ActionWithIcon
	Target string `yaml:"target"`
}

func (a *SwapPage) Run(pressed bool) error {
	if !pressed {
		return nil
	}
	return a.Controller.SwapPage(a.Target)
}

func (a *SwapPage) Name() string {
	return "swap_page"
}
