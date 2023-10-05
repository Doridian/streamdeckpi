package homeassistant

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/Doridian/go-haws"
	"github.com/Doridian/go-streamdeck"
	"github.com/Doridian/streamdeckpi/agent/internal/action"
	"github.com/Doridian/streamdeckpi/agent/internal/controller"
	"gopkg.in/yaml.v3"
)

type haStringAction struct {
	haEntityActionBase

	Icon string `yaml:"icon"`

	state    string
	doRender bool
}

func (a *haStringAction) New() action.Action {
	return &haStringAction{}
}

func (a *haStringAction) OnState(entityID string, state haws.State) error {
	a.state = state.State
	a.doRender = true
	return nil
}

func (a *haStringAction) ApplyConfig(config *yaml.Node, imageHelper controller.ImageHelper, ctrl controller.Controller) error {
	err := a.haEntityActionBase.ApplyConfig(config, imageHelper, ctrl)
	if err != nil {
		return err
	}

	err = config.Decode(a)
	if err != nil {
		return err
	}

	a.instance.RegisterStateReceiver(a, a.Entity)

	return nil
}

func (a *haStringAction) Name() string {
	return "homeassistant_string"
}

func (a *haStringAction) Render(force bool) (*streamdeck.ImageData, error) {
	if !force && !a.doRender {
		return nil, nil
	}

	var err error
	var convImg *streamdeck.ImageData

	baseImage, err := a.ImageHelper.LoadNoConvert(a.Icon)
	if err != nil {
		return nil, err
	}

	img := image.NewRGBA(baseImage.Bounds())

	draw.Draw(img, img.Rect, baseImage, image.Point{}, draw.Src)
	drawCenteredText(a.Controller, img, color.RGBA{255, 255, 255, 255}, 48, 48, a.state)
	convImg, err = a.ImageHelper.Convert(img)

	if err == nil {
		a.doRender = false
	}

	return convImg, err
}
