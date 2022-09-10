package impl

import (
	"errors"
	"log"
	"time"

	"github.com/Doridian/go-streamdeck"
	"github.com/Doridian/streamdeckpi/agent/action"
)

func (c *controllerImpl) renderAction(actionObj action.Action, force bool) (*streamdeck.ImageData, error) {
	if actionObj == nil {
		return nil, nil
	}

	return actionObj.Render(force)
}

func (c *controllerImpl) render(force bool) (hadErrors bool) {
	currentPage := c.pageTop
	if currentPage != c.lastRenderedPage {
		force = true
	}
	c.lastRenderedPage = currentPage

	var img *streamdeck.ImageData
	var err error

	for i, actionObj := range c.pageTop.actions {
		img, err = c.renderAction(actionObj, force)
		if err != nil {
			log.Printf("Error rendering action: %v", err)
			hadErrors = true
			img = nil
		}

		if force && img == nil {
			img = c.imageLoader.GetBlankImage()
		}

		if img == nil {
			continue
		}

		err = c.dev.SetConvertedImage(uint8(i), img)
		if err != nil {
			log.Printf("Error setting image: %v", err)
			hadErrors = true
		}
	}

	return
}

func (c *controllerImpl) renderLoop() {
	defer c.runWait.Done()

	frameWait := time.Duration(16) * time.Millisecond
	errorWait := time.Duration(1) * time.Second

	hadErrors := false
	for c.running {
		hadErrors = c.render(hadErrors)
		if hadErrors {
			time.Sleep(errorWait)
		}
		time.Sleep(frameWait)
	}

	c.stopError(errors.New("unexpectedly reached end of render loop"))
}
