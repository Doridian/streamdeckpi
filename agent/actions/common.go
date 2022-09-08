package actions

import (
	"github.com/Doridian/streamdeck"
	"github.com/Doridian/streamdeckpi/agent/interfaces"
)

type actionWithIcon struct {
	imageLoader interfaces.ImageLoader
	Icon        string `yaml:"icon"`
}

func (a *actionWithIcon) ApplyConfig(imageLoader interfaces.ImageLoader) error {
	a.imageLoader = imageLoader
	return nil
}

func (a *actionWithIcon) Render(force bool) (*streamdeck.ImageData, error) {
	if force {
		return a.imageLoader.Load(a.Icon)
	}
	return nil, nil
}

func (a *actionWithIcon) GetConfigRef() interface{} {
	return a
}
