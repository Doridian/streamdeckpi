package page

import "github.com/Doridian/streamdeckpi/agent/action"

type PushPage struct {
	action.ActionWithIcon
	Target string `yaml:"target"`
}

func (a *PushPage) Run(pressed bool) error {
	if !pressed {
		return nil
	}
	return a.Controller.PushPage(a.Target)
}

func (a *PushPage) Name() string {
	return "push_page"
}
