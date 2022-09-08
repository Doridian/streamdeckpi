package controller

import (
	"errors"
	"log"
	"time"

	"github.com/Doridian/streamdeck"
	"github.com/Doridian/streamdeckpi/agent/actions"
)

func (c *controller) renderAction(action actions.Action, force bool) (*streamdeck.ImageData, error) {
	if action == nil {
		return nil, nil
	}

	return action.Render(force)
}

func (c *controller) render() (hadErrors bool) {
	currentPage := c.pageTop
	pageSwapped := currentPage != c.lastRenderedPage
	c.lastRenderedPage = currentPage

	var img *streamdeck.ImageData
	var err error

	for i, action := range c.pageTop.actions {
		img, err = c.renderAction(action, pageSwapped)
		if err != nil {
			log.Printf("Error rendering action: %v", err)
			hadErrors = true
			img = nil
		}

		if pageSwapped && img == nil {
			img = c.blankImage
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

func (c *controller) renderLoop() {
	defer c.runWait.Done()

	frameWait := time.Duration(16) * time.Millisecond
	errorWait := time.Duration(1) * time.Second

	var hadErrors bool
	for c.running {
		hadErrors = c.render()
		if hadErrors {
			time.Sleep(errorWait)
		}
		time.Sleep(frameWait)
	}

	c.stopError(errors.New("unexpectedly reached end of render loop"))
}
