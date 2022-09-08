package actions

import (
	"github.com/Doridian/streamdeck"
	"github.com/Doridian/streamdeckpi/agent/interfaces"
)

type ActionBase struct {
	ImageLoader interfaces.ImageLoader
	Controller  interfaces.Controller
}

func (a *ActionBase) ApplyConfig(imageLoader interfaces.ImageLoader, controller interfaces.Controller) error {
	a.ImageLoader = imageLoader
	a.Controller = controller
	return nil
}

func (a *ActionBase) GetConfigRef() interface{} {
	return a
}

type ActionWithIcon struct {
	ActionBase
	Icon string `yaml:"icon"`
}

func (a *ActionWithIcon) Render(force bool) (*streamdeck.ImageData, error) {
	if force {
		return a.ImageLoader.Load(a.Icon)
	}
	return nil, nil
}
