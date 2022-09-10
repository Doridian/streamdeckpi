package controller

import "github.com/Doridian/go-streamdeck"

type ImageLoader interface {
	Load(path string) (*streamdeck.ImageData, error)
	GetBlankImage() *streamdeck.ImageData
}
