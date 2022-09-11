package homeassistant

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/Doridian/go-haws"
	"github.com/Doridian/go-streamdeck"
	"github.com/Doridian/streamdeckpi/agent/action"
	"github.com/Doridian/streamdeckpi/agent/controller"
	"gopkg.in/yaml.v3"
)

type haLightAction struct {
	haEntityActionBase

	baseImage  image.Image
	renderIcon *streamdeck.ImageData
	doRender   bool
	lightColor color.Color
}

func (a *haLightAction) New() action.Action {
	return &haLightAction{}
}

func convColorElement(elem interface{}, brightness float64) uint8 {
	col := elem.(float64)
	col *= (brightness / 255.0)
	return uint8(col)
}

func (a *haLightAction) OnState(entityID string, state haws.State) error {
	if state.State == "off" {
		a.lightColor = color.Black
	} else {
		lightColorRGB := state.Attributes["rgb_color"].([]interface{})
		brightness := state.Attributes["brightness"].(float64)

		a.lightColor = color.NRGBA{
			R: convColorElement(lightColorRGB[0], brightness),
			G: convColorElement(lightColorRGB[1], brightness),
			B: convColorElement(lightColorRGB[2], brightness),
			A: 255,
		}
	}

	a.doRender = true
	return nil
}

func (a *haLightAction) ApplyConfig(config *yaml.Node, imageHelper controller.ImageHelper, ctrl controller.Controller) error {
	err := a.haEntityActionBase.ApplyConfig(config, imageHelper, ctrl)
	if err != nil {
		return err
	}

	err = config.Decode(a)
	if err != nil {
		return err
	}

	if a.Domain == "" {
		a.Domain = "light"
	}

	baseImage, err := imageHelper.LoadNoConvert(a.Icon)
	if err != nil {
		return err
	}
	a.baseImage = baseImage
	a.renderIcon, err = imageHelper.Convert(a.baseImage)
	if err != nil {
		return err
	}

	a.instance.RegisterStateReceiver(a, a.Entity)

	return nil
}

func (a *haLightAction) Name() string {
	return "homeassistant_light"
}

func (a *haLightAction) Render(force bool) (*streamdeck.ImageData, error) {
	if !force && !a.doRender {
		return nil, nil
	}

	img := image.NewRGBA(a.baseImage.Bounds())

	draw.Draw(img, img.Rect, image.NewUniform(a.lightColor), image.Point{}, draw.Src)
	draw.Draw(img, img.Rect, a.baseImage, image.Point{}, draw.Over)

	convImg, err := a.ImageHelper.Convert(img)
	if err == nil {
		a.doRender = false
	}

	return convImg, err
}
