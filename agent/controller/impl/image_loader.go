package impl

import (
	"image"
	"path"
	"sync"

	"github.com/Doridian/streamdeck"
	"github.com/Doridian/streamdeckpi/agent/controller"
)

type imageLoader struct {
	path       string
	controller *controllerImpl

	imageCache     map[string]*streamdeck.ImageData
	imageCacheLock sync.RWMutex
}

func newImageLoader(controller *controllerImpl, page *page) controller.ImageLoader {
	return &imageLoader{
		path:       path.Dir(page.path),
		controller: controller,
	}
}

func (l *imageLoader) Load(path string) (*streamdeck.ImageData, error) {
	path, err := l.controller.cleanPath(path)
	if err != nil {
		return nil, err
	}

	l.imageCacheLock.RLock()
	img, ok := l.imageCache[path]
	l.imageCacheLock.RUnlock()
	if ok {
		return img, nil
	}

	reader, err := l.controller.resolveFile(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	goImage, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	convImg, err := l.controller.dev.ConvertImage(goImage)
	if err != nil {
		return nil, err
	}

	l.imageCacheLock.Lock()
	l.imageCache[path] = convImg
	l.imageCacheLock.Unlock()

	return convImg, nil
}
