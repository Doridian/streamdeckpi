package impl

import (
	"fmt"
	"image"
	"path"
	"sync"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

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
		imageCache: make(map[string]*streamdeck.ImageData),
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
		return nil, fmt.Errorf("error resolving image: %w", err)
	}
	defer reader.Close()

	goImage, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}

	convImg, err := l.controller.dev.ConvertImage(goImage)
	if err != nil {
		return nil, fmt.Errorf("error converting image: %w", err)
	}

	l.imageCacheLock.Lock()
	l.imageCache[path] = convImg
	l.imageCacheLock.Unlock()

	return convImg, nil
}
