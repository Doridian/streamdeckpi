package controller

import (
	"image"
	"path"

	"github.com/Doridian/streamdeck"
	"github.com/Doridian/streamdeckpi/agent/interfaces"
)

type imageLoader struct {
	path       string
	controller *controller
}

func newImageLoader(controller *controller, page *page) interfaces.ImageLoader {
	return &imageLoader{
		path:       path.Dir(page.path),
		controller: controller,
	}
}

func (l *imageLoader) Load(path string) (*streamdeck.ImageData, error) {
	reader, err := l.controller.resolveFile(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return l.controller.dev.ConvertImage(img)
}
