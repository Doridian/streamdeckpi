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

type haStringConditionOverride struct {
	Condition haCondition `yaml:"condition"`
	Icon      string      `yaml:"icon"`

	Color         []uint8 `yaml:"color"`
	X             *int    `yaml:"x"`
	Y             *int    `yaml:"y"`
	Size          float64 `yaml:"size"`
	Font          string  `yaml:"font"`
	Align         string  `yaml:"align"`
	VerticalAlign string  `yaml:"vertical-align"`
}

type haStringAction struct {
	haEntityActionBase

	Icon          string                       `yaml:"icon"`
	Conditions    []*haStringConditionOverride `yaml:"conditions"`
	Color         []uint8                      `yaml:"color"`
	X             int                          `yaml:"x"`
	Y             int                          `yaml:"y"`
	Size          float64                      `yaml:"size"`
	Font          string                       `yaml:"font"`
	Align         string                       `yaml:"align"`
	VerticalAlign string                       `yaml:"vertical-align"`

	useFont          string
	useColor         color.RGBA
	useX             int
	useY             int
	useSize          float64
	useIcon          string
	useAlign         FontAlign
	useVerticalAlign FontAlign
	state            string
	doRender         bool
}

func (a *haStringAction) New() action.Action {
	return &haStringAction{}
}

func (a *haStringAction) OnState(entityID string, state haws.State) error {
	var currentMatch *haStringConditionOverride
	for _, cond := range a.Conditions {
		match, err := cond.Condition.Evaluate(&state)
		if err != nil {
			return err
		}

		if match {
			currentMatch = cond
			break
		}
	}

	newUseIcon := a.Icon
	newUseColor := a.Color
	newUseX := a.X
	newUseY := a.Y
	newUseSize := a.Size
	newUseFont := a.Font
	newUseAlign := a.Align
	newUseVerticalAlign := a.VerticalAlign

	if currentMatch != nil {
		if currentMatch.Icon != "" {
			newUseIcon = currentMatch.Icon
		}
		if currentMatch.Color != nil {
			newUseColor = currentMatch.Color
		}
		if currentMatch.X != nil {
			newUseX = *currentMatch.X
		}
		if currentMatch.Y != nil {
			newUseY = *currentMatch.Y
		}
		if currentMatch.Size != 0 {
			newUseSize = currentMatch.Size
		}
		if currentMatch.Font != "" {
			newUseFont = currentMatch.Font
		}
		if currentMatch.Align != "" {
			newUseAlign = currentMatch.Align
		}
		if currentMatch.VerticalAlign != "" {
			newUseVerticalAlign = currentMatch.VerticalAlign
		}
	}

	if newUseFont == "" {
		newUseFont = "font.ttf"
	}
	if newUseSize <= 0 {
		newUseSize = 48
	}
	if newUseColor == nil {
		newUseColor = []uint8{255, 255, 255}
	}
	if newUseAlign == "" {
		newUseAlign = "center"
	}
	if newUseVerticalAlign == "" {
		newUseVerticalAlign = "middle"
	}

	a.useColor = color.RGBA{newUseColor[0], newUseColor[1], newUseColor[2], 255}
	a.useX = newUseX
	a.useY = newUseY
	a.useIcon = newUseIcon
	a.useSize = newUseSize
	a.useFont = newUseFont
	switch newUseAlign {
	case "left":
		a.useAlign = FontAlignLeft
	case "center":
		a.useAlign = FontAlignCenter
	case "right":
		a.useAlign = FontAlignRight
	}
	switch newUseVerticalAlign {
	case "top":
		a.useVerticalAlign = FontAlignTop
	case "middle":
		a.useVerticalAlign = FontAlignMiddle
	case "bottom":
		a.useVerticalAlign = FontAlignBottom
	}

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

	baseImage, err := a.ImageHelper.LoadNoConvert(a.useIcon)
	if err != nil {
		return nil, err
	}

	img := image.NewRGBA(baseImage.Bounds())

	draw.Draw(img, img.Rect, baseImage, image.Point{}, draw.Src)
	drawCenteredText(a.Controller, img, a.useFont, a.useColor, a.useX, a.useY, float64(a.useSize), a.useAlign, a.useVerticalAlign, a.state)
	convImg, err = a.ImageHelper.Convert(img)

	if err == nil {
		a.doRender = false
	}

	return convImg, err
}
