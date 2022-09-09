package page

import "github.com/Doridian/streamdeckpi/agent/actions"

type PopPage struct {
	actions.ActionWithIcon
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
