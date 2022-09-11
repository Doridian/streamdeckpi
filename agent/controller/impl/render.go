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

func (c *controllerImpl) render() {
	currentPage := c.pageTop
	if currentPage != c.lastRenderedPage {
		c.renderOkState = 0
	}
	c.lastRenderedPage = currentPage

	var img *streamdeck.ImageData

	var renderErr error
	var setErr error

	for i, actionObj := range c.pageTop.actions {
		renderOkMask := uint64(0b1) << i

		force := c.renderOkState&renderOkMask == 0

		img, renderErr = c.renderAction(actionObj, force)
		if renderErr != nil {
			log.Printf("Error rendering action: %v", renderErr)
			img = nil
		}

		if force && img == nil {
			img = c.imageHelper.GetBlankImage()
		}

		if img == nil {
			setErr = nil
		} else {
			setErr = c.dev.SetConvertedImage(uint8(i), img)
			if setErr != nil {
				log.Printf("Error setting image: %v", setErr)
			}
		}

		if renderErr == nil && setErr == nil {
			c.renderOkState |= renderOkMask
		}
	}
}

func (c *controllerImpl) renderLoop() {
	defer c.runWait.Done()

	frameWait := time.Duration(16) * time.Millisecond

	for c.running {
		c.render()
		time.Sleep(frameWait)
	}

	c.stopError(errors.New("unexpectedly reached end of render loop"))
}
