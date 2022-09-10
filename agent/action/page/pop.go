package page

import "github.com/Doridian/streamdeckpi/agent/action"

type popPage struct {
	action.ActionWithIcon
}

func (a *popPage) New() action.Action {
	return &popPage{}
}

func (a *popPage) Run(pressed bool) error {
	if !pressed {
		return nil
	}
	return a.Controller.PopPage()
}

func (a *popPage) Name() string {
	return "pop_page"
}
