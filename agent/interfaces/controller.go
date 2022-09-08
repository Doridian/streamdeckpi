package interfaces

import "github.com/Doridian/streamdeck"

type Controller interface {
	SwapPage(pageFile string) error
	PushPage(pageFile string) error
	PopPage() error

	Start() error
	Stop() error
	Wait() error

	GetBlankImage() *streamdeck.ImageData
}
