package interfaces

import "github.com/Doridian/streamdeck"

type ImageLoader interface {
	Load(path string) (*streamdeck.ImageData, error)
}
