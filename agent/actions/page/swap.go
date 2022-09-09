package page

import "github.com/Doridian/streamdeckpi/agent/actions"

type SwapPage struct {
	actions.ActionWithIcon
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
