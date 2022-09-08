package controller

import (
	"errors"
	"image"
	"sync"

	"github.com/Doridian/streamdeck"
	"github.com/Doridian/streamdeckpi/agent/interfaces"
)

var errStopNone = errors.New("no error")

type controller struct {
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
}

func NewController(dev *streamdeck.Device) (interfaces.Controller, error) {
	img := image.NewRGBA(image.Rect(0, 0, int(dev.Pixels), int(dev.Pixels)))

	convImg, err := dev.ConvertImage(img)
	if err != nil {
		return nil, err
	}

	return &controller{
		pageStack:  make([]*page, 0),
		pageTop:    nil,
		dev:        dev,
		running:    false,
		blankImage: convImg,
	}, nil
}

func (c *controller) GetBlankImage() *streamdeck.ImageData {
	return c.blankImage
}

func (c *controller) Start() error {
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

func (c *controller) Reset() error {
	c.pageWait.Lock()
	defer c.pageWait.Unlock()
	c.pageCacheLock.Lock()
	defer c.pageCacheLock.Unlock()

	c.pageCache = make(map[string]*page)
	c.pageStack = make([]*page, 0)
	c.pageTop = nil
	return c.PushPage("./default.yml")
}

func (c *controller) Stop() error {
	return c.stopSync(errStopNone)
}

func (c *controller) Wait() error {
	c.runWait.Wait()

	err := c.runError
	if err == errStopNone {
		return nil
	}
	return err
}

func (c *controller) stopError(err error) {
	go c.stopSync(err)
}

func (c *controller) stopSync(err error) error {
	c.runControlWait.Lock()
	defer c.runControlWait.Unlock()

	if c.runError == nil && err != nil {
		c.runError = err
	}

	c.running = false
	c.runWait.Wait()
	return c.dev.Close()
}
