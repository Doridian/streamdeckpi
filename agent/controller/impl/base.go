package impl

import (
	"errors"
	"sync"

	"github.com/Doridian/streamdeck"
	"github.com/Doridian/streamdeckpi/agent/controller"
)

var errStopNone = errors.New("no error")

type controllerImpl struct {
	pageStack []*page
	pageTop   *page
	pageWait  sync.Mutex

	lastRenderedPage *page
	blankImage       *streamdeck.ImageData

	running        bool
	runWait        sync.WaitGroup
	runError       error
	runControlWait sync.Mutex

	dev *streamdeck.Device

	pageCache     map[string]*page
	pageCacheLock sync.Mutex

	imageLoader controller.ImageLoader
}

func NewController(dev *streamdeck.Device) (controller.Controller, error) {
	ctrl := &controllerImpl{
		pageStack: make([]*page, 0),
		pageTop:   nil,
		dev:       dev,
		running:   false,
	}

	var err error
	ctrl.imageLoader, err = newImageLoader(ctrl)

	return ctrl, err
}

func (c *controllerImpl) Start() error {
	c.runControlWait.Lock()
	defer c.runControlWait.Unlock()

	if c.running {
		return errors.New("already running")
	}

	err := c.Reset()
	if err != nil {
		return err
	}

	err = c.dev.Open()
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

func (c *controllerImpl) Reset() error {
	c.pageCacheLock.Lock()
	c.pageCache = make(map[string]*page)
	c.pageCacheLock.Unlock()

	pageObj, err := c.resolvePage("default.yml")
	if err != nil {
		return err
	}

	c.pageWait.Lock()
	c.pageStack = []*page{pageObj}
	c.pageTop = pageObj
	c.pageWait.Unlock()

	return nil
}

func (c *controllerImpl) Stop() error {
	return c.stopSync(errStopNone)
}

func (c *controllerImpl) Wait() error {
	c.runWait.Wait()

	err := c.runError
	if err == errStopNone {
		return nil
	}
	return err
}

func (c *controllerImpl) stopError(err error) {
	go c.stopSync(err)
}

func (c *controllerImpl) stopSync(err error) error {
	c.runControlWait.Lock()
	defer c.runControlWait.Unlock()

	if c.runError == nil && err != nil {
		c.runError = err
	}

	c.running = false
	c.runWait.Wait()
	return c.dev.Close()
}
