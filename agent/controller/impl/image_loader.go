package impl

import (
	"fmt"
	"image"
	"sync"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/Doridian/go-streamdeck"
	"github.com/Doridian/streamdeckpi/agent/controller"
)

type imageLoader struct {
	controller *controllerImpl

	imageCache     map[string]*streamdeck.ImageData
	imageCacheLock sync.RWMutex
	blankImage     *streamdeck.ImageData
}

func newImageLoader(ctrl *controllerImpl) (controller.ImageLoader, error) {
	img := image.NewRGBA(image.Rect(0, 0, int(ctrl.dev.Pixels), int(ctrl.dev.Pixels)))

	convImg, err := ctrl.dev.ConvertImage(img)
	if err != nil {
		return nil, err
	}

	return &imageLoader{
		controller: ctrl,
		blankImage: convImg,
		imageCache: make(map[string]*streamdeck.ImageData),
	}, nil
}

func (l *imageLoader) GetBlankImage() *streamdeck.ImageData {
	return l.blankImage
}

func (l *imageLoader) LoadNoConvert(pathSub string) (image.Image, error) {
	reader, err := l.controller.ResolveFile(pathSub)
	if err != nil {
		return nil, fmt.Errorf("error resolving image: %w", err)
	}
	defer reader.Close()

	goImage, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}

	return goImage, nil
}

func (l *imageLoader) Convert(img image.Image) (*streamdeck.ImageData, error) {
	return l.controller.dev.ConvertImage(img)
}

func (l *imageLoader) Load(pathSub string) (*streamdeck.ImageData, error) {
	pathSub, err := l.controller.CleanPath(pathSub)
	if err != nil {
		return nil, err
	}

	l.imageCacheLock.RLock()
	img, ok := l.imageCache[pathSub]
	l.imageCacheLock.RUnlock()
	if ok {
		return img, nil
	}

	goImage, err := l.LoadNoConvert(pathSub)
	if err != nil {
		return nil, err
	}

	convImg, err := l.Convert(goImage)
	if err != nil {
		return nil, fmt.Errorf("error converting image: %w", err)
	}

	l.imageCacheLock.Lock()
	l.imageCache[pathSub] = convImg
	l.imageCacheLock.Unlock()

	return convImg, nil
}
