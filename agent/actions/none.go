package actions

type None struct {
	ActionWithIcon
}

func (a *None) Run(pressed bool) error {
	return nil
}

func (a *None) Name() string {
	return "none"
}
