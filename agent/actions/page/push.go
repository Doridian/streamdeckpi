package page

import "github.com/Doridian/streamdeckpi/agent/actions"

type PushPage struct {
	actions.ActionWithIcon
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
