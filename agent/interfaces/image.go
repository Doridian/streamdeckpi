package interfaces

import "image"

type ImageLoader = func(path string) (image.Image, error)
