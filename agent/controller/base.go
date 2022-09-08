package controller

import (
	"errors"
	"image"
	"io"
	"sync"

	"github.com/Doridian/streamdeck"
	"github.com/Doridian/streamdeckpi/agent/interfaces"
)

var errStopNone = errors.New("no error")

type controller struct {
	pageStack []*page
	pageTop   *page

	lastRenderedPage *page
	blankImage       image.Image

	running        bool
	runWait        sync.WaitGroup
	runError       error
	runControlWait sync.Mutex

	dev *streamdeck.Device
}

func NewController(dev *streamdeck.Device) (interfaces.Controller, error) {
	res := &controller{
		pageStack:  make([]*page, 0),
		dev:        dev,
		running:    false,
		blankImage: image.NewRGBA(image.Rect(0, 0, int(dev.Pixels), int(dev.Pixels))),
	}

	err := res.PushPage("default.yml")
	if err != nil {
		err = res.PushPage("/:embed/default.yml")
	}

	return res, err
}

func (c *controller) Start() error {
	c.runControlWait.Lock()
	defer c.runControlWait.Unlock()

	err := c.dev.Open()
	if err != nil {
		return err
	}

	c.runError = nil
	c.running = true

	c.runWait.Add(1)
	go c.buttonLoop()

	c.runWait.Add(1)
	go c.renderLoop()

	return nil
}

func (c *controller) Stop() error {
	c.runControlWait.Lock()
	defer c.runControlWait.Unlock()

	return c.stopSync(errStopNone)
}

func (c *controller) Wait() error {
	c.runWait.Wait()
	if c.runError == errStopNone {
		return nil
	}
	return c.runError
}

func (c *controller) stopError(err error) {
	go c.stopSync(err)
}

func (c *controller) stopSync(err error) error {
	c.runControlWait.Lock()
	defer c.runControlWait.Unlock()

	if c.runError == nil {
		c.runError = err
	}

	c.running = false
	c.Wait()
	return c.dev.Close()
}

func (c *controller) resolveFile(file string) (io.ReadCloser, string, error) {
	return nil, "", errors.New("could not find file")
}
