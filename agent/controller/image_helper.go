package controller

import (
	"image"

	"github.com/Doridian/go-streamdeck"
)

type ImageHelper interface {
	Load(path string) (*streamdeck.ImageData, error)

	LoadNoConvert(pathSub string) (image.Image, error)
	Convert(img image.Image) (*streamdeck.ImageData, error)

	GetBlankImage() *streamdeck.ImageData
}
