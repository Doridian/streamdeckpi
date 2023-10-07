package homeassistant

import (
	"image"
	"image/color"
	"image/draw"
	"strings"

	"github.com/Doridian/go-haws"
	"github.com/Doridian/go-streamdeck"
	"github.com/Doridian/streamdeckpi/agent/internal/action"
	"github.com/Doridian/streamdeckpi/agent/internal/controller"
	"gopkg.in/yaml.v3"
)

type haStringConditionOverride struct {
	Condition haCondition `yaml:"condition"`
	Icon      string      `yaml:"icon"`

	Texts []*haStringActionText `yaml:"texts"`
}

type haStringActionText struct {
	Color         []uint8  `yaml:"color"`
	X             *int     `yaml:"x"`
	Y             *int     `yaml:"y"`
	Size          *float64 `yaml:"size"`
	Font          *string  `yaml:"font"`
	Align         *string  `yaml:"align"`
	VerticalAlign *string  `yaml:"vertical-align"`
	Text          *string  `yaml:"text"`

	color         color.RGBA
	align         FontAlign
	verticalAlign FontAlign
	text          string
}

type haStringAction struct {
	haEntityActionBase

	Icon       string                       `yaml:"icon"`
	Conditions []*haStringConditionOverride `yaml:"conditions"`
	Texts      []haStringActionText         `yaml:"texts"`

	useTexts []haStringActionText
	useIcon  string

	state    string
	doRender bool
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
	newUseTexts := make([]haStringActionText, 0, len(a.Texts))

	if currentMatch != nil && currentMatch.Icon != "" {
		newUseIcon = currentMatch.Icon
	}

	for i, text := range a.Texts {
		if currentMatch != nil {
			currentMatchTexts := currentMatch.Texts[i]
			if currentMatchTexts != nil {
				if currentMatchTexts.Text != nil {
					text.Text = currentMatchTexts.Text
				}
				if currentMatchTexts.Color != nil {
					text.Color = currentMatchTexts.Color
				}
				if currentMatchTexts.X != nil {
					text.X = currentMatchTexts.X
				}
				if currentMatchTexts.Y != nil {
					text.Y = currentMatchTexts.Y
				}
				if currentMatchTexts.Size != nil {
					text.Size = currentMatchTexts.Size
				}
				if currentMatchTexts.Font != nil {
					text.Font = currentMatchTexts.Font
				}
				if currentMatchTexts.Align != nil {
					text.Align = currentMatchTexts.Align
				}
				if currentMatchTexts.VerticalAlign != nil {
					text.VerticalAlign = currentMatchTexts.VerticalAlign
				}
			}
		}

		switch *text.Align {
		case "left":
			text.align = FontAlignLeft
		case "right":
			text.align = FontAlignRight
		case "center":
			fallthrough
		default:
			text.align = FontAlignCenter
		}
		switch *text.VerticalAlign {
		case "top":
			text.verticalAlign = FontAlignTop
		case "bottom":
			text.verticalAlign = FontAlignBottom
		case "middle":
			fallthrough
		default:
			text.verticalAlign = FontAlignMiddle

		}

		text.color = color.RGBA{text.Color[0], text.Color[1], text.Color[2], text.Color[3]}
		text.text = strings.ReplaceAll(*text.Text, "$STATE", state.State)

		newUseTexts = append(newUseTexts, text)
	}

	a.useTexts = newUseTexts
	a.useIcon = newUseIcon

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
	for _, text := range a.useTexts {
		err = drawCenteredText(a.Controller, img, *text.Font, text.color, *text.X, *text.Y, *text.Size, text.align, text.verticalAlign, text.text)
		if err != nil {
			return nil, err
		}
	}

	convImg, err = a.ImageHelper.Convert(img)
	if err != nil {
		return nil, err
	}

	a.doRender = false
	return convImg, nil
}
