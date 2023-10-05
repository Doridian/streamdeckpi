package impl

import (
	"fmt"
	"image"
	"sync"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/Doridian/go-streamdeck"
	"github.com/Doridian/streamdeckpi/agent/internal/controller"
)

type imageHelper struct {
	controller *controllerImpl

	imageCache     map[string]*streamdeck.ImageData
	imageCacheLock *sync.RWMutex

	rawImageCache     map[string]image.Image
	rawImageCacheLock *sync.RWMutex

	blankImage  *streamdeck.ImageData
	imageBounds image.Rectangle
}

func newImageHelper(ctrl *controllerImpl) (controller.ImageHelper, error) {
	bounds := image.Rect(0, 0, int(ctrl.dev.Pixels), int(ctrl.dev.Pixels))
	img := image.NewRGBA(bounds)

	convImg, err := ctrl.dev.ConvertImage(img)
	if err != nil {
		return nil, err
	}

	return &imageHelper{
		controller: ctrl,

		blankImage:  convImg,
		imageBounds: bounds,

		imageCache:     make(map[string]*streamdeck.ImageData),
		imageCacheLock: &sync.RWMutex{},

		rawImageCache:     make(map[string]image.Image),
		rawImageCacheLock: &sync.RWMutex{},
	}, nil
}

func (l *imageHelper) GetBlankImage() *streamdeck.ImageData {
	return l.blankImage
}

func (l *imageHelper) GetImageBounds() image.Rectangle {
	return l.imageBounds
}

func (l *imageHelper) LoadNoConvert(pathSub string) (image.Image, error) {
	pathSub, err := l.controller.CleanPath(pathSub)
	if err != nil {
		return nil, err
	}

	l.rawImageCacheLock.RLock()
	img, ok := l.rawImageCache[pathSub]
	l.rawImageCacheLock.RUnlock()
	if ok {
		return img, nil
	}

	// Re-try fetching here, just in case another run just created the cache entry
	// as we could't hold a lock for a little while
	l.rawImageCacheLock.Lock()
	defer l.rawImageCacheLock.Unlock()

	img, ok = l.rawImageCache[pathSub]
	if ok {
		return img, nil
	}

	img, err = l.loadNoConvert(pathSub)
	if err != nil {
		return nil, err
	}

	l.rawImageCache[pathSub] = img

	return img, nil
}

func (l *imageHelper) loadNoConvert(pathSub string) (image.Image, error) {
	reader, err := l.controller.ResolveFile(pathSub)
	if err != nil {
		return nil, fmt.Errorf("error resolving image: %w", err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}

	return img, nil
}

func (l *imageHelper) Convert(img image.Image) (*streamdeck.ImageData, error) {
	return l.controller.dev.ConvertImage(img)
}

func (l *imageHelper) Load(pathSub string) (*streamdeck.ImageData, error) {
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

	l.imageCacheLock.Lock()
	defer l.imageCacheLock.Unlock()

	// Re-try fetching here, just in case another run just created the cache entry
	// as we could't hold a lock for a little while
	img, ok = l.imageCache[pathSub]
	if ok {
		return img, nil
	}

	goImage, err := l.loadNoConvert(pathSub)
	if err != nil {
		return nil, err
	}

	convImg, err := l.Convert(goImage)
	if err != nil {
		return nil, fmt.Errorf("error converting image: %w", err)
	}

	l.imageCache[pathSub] = convImg

	return convImg, nil
}
