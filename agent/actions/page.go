package actions

type SwapPage struct {
	ActionWithIcon
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

type PopPage struct {
	ActionWithIcon
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

type PushPage struct {
	ActionWithIcon
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
