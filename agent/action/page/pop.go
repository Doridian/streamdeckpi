package page

import "github.com/Doridian/streamdeckpi/agent/action"

type PopPage struct {
	action.ActionWithIcon
}

func (a *PopPage) New() action.Action {
	return &PopPage{}
}

func (a *PopPage) Run(pressed bool) error {
	if !pressed {
		return nil
	}
	return a.Controller.PopPage()
}

func (a *PopPage) Name() string {
	return "pop_page"
}
