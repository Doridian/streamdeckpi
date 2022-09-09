package controller

import "github.com/Doridian/streamdeck"

type Controller interface {
	SwapPage(pageFile string) error
	PushPage(pageFile string) error
	PopPage() error

	Start() error
	Stop() error
	Wait() error
	Reset() error

	GetBlankImage() *streamdeck.ImageData
}
