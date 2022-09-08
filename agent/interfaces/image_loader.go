package interfaces

import "image"

type ImageLoader interface {
	Load(path string) (image.Image, error)
}
