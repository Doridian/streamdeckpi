package impl

import (
	"image"
	"path"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/Doridian/go-streamdeck"
	"github.com/Doridian/streamdeckpi/agent/controller"
)

type pageImageLoader struct {
	path   string
	parent controller.ImageLoader
}

func newPageImageLoader(parent controller.ImageLoader, pageObj *page) controller.ImageLoader {
	return &pageImageLoader{
		path:   path.Dir(pageObj.path),
		parent: parent,
	}
}

func (l *pageImageLoader) Load(pathSub string) (*streamdeck.ImageData, error) {
	return l.parent.Load(path.Join(l.path, pathSub))
}

func (l *pageImageLoader) GetBlankImage() *streamdeck.ImageData {
	return l.parent.GetBlankImage()
}

func (l *pageImageLoader) LoadNoConvert(pathSub string) (image.Image, error) {
	return l.parent.LoadNoConvert(path.Join(l.path, pathSub))
}

func (l *pageImageLoader) Convert(img image.Image) (*streamdeck.ImageData, error) {
	return l.parent.Convert(img)
}
