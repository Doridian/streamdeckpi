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

type pageImageHelper struct {
	path   string
	parent controller.ImageHelper
}

func newPageImageHelper(parent controller.ImageHelper, pageObj *page) controller.ImageHelper {
	return &pageImageHelper{
		path:   path.Dir(pageObj.path),
		parent: parent,
	}
}

func (l *pageImageHelper) Load(pathSub string) (*streamdeck.ImageData, error) {
	return l.parent.Load(path.Join(l.path, pathSub))
}

func (l *pageImageHelper) GetBlankImage() *streamdeck.ImageData {
	return l.parent.GetBlankImage()
}

func (l *pageImageHelper) GetImageBounds() image.Rectangle {
	return l.parent.GetImageBounds()
}

func (l *pageImageHelper) LoadNoConvert(pathSub string) (image.Image, error) {
	return l.parent.LoadNoConvert(path.Join(l.path, pathSub))
}

func (l *pageImageHelper) Convert(img image.Image) (*streamdeck.ImageData, error) {
	return l.parent.Convert(img)
}
