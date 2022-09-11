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

	OnIcon  string `yaml:"on_icon"`
	OffIcon string `yaml:"off_icon"`

	baseImageOn image.Image
	imageOff    *streamdeck.ImageData

	doRender bool

	lightColor color.Color
	lightOn    bool
}

func (a *haLightAction) New() action.Action {
	return &haLightAction{}
}

func convColorElement(elem interface{}, brightness float64) uint8 {
	col, ok := coerceNumber(elem)
	if !ok {
		return 255
	}
	col *= (brightness / 255.0)
	return uint8(col)
}

func (a *haLightAction) OnState(entityID string, state haws.State) error {
	if state.State == "off" {
		a.lightOn = false

		a.lightColor = color.Black
	} else {
		a.lightOn = true

		lightColorRGB, ok := state.Attributes["rgb_color"].([]interface{})
		if !ok {
			lightColorRGB = []interface{}{
				255,
				255,
				255,
			}
		}

		brightness, ok := coerceNumber(state.Attributes["brightness"])
		if !ok {
			brightness = 255
		}

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

	a.baseImageOn, err = imageHelper.LoadNoConvert(a.OnIcon)
	if err != nil {
		return err
	}

	a.imageOff, err = imageHelper.Load(a.OffIcon)
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

	var err error
	var convImg *streamdeck.ImageData
	if a.lightOn {
		img := image.NewRGBA(a.baseImageOn.Bounds())
		draw.Draw(img, img.Rect, image.NewUniform(a.lightColor), image.Point{}, draw.Src)
		draw.Draw(img, img.Rect, a.baseImageOn, image.Point{}, draw.Over)
		convImg, err = a.ImageHelper.Convert(img)
	} else {
		convImg = a.imageOff
	}

	if err == nil {
		a.doRender = false
	}
	return convImg, err
}
