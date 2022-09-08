package controller

import (
	"errors"
	"io"
	"time"

	"github.com/Doridian/streamdeckpi/agent/actions"
	"gopkg.in/yaml.v3"
)

type page struct {
	path string

	Timeout time.Duration `yaml:"timeout"`

	Actions []actions.Action `yaml:"actions"`
}

func (c *controller) resolvePage(pageFile string) (*page, error) {
	reader, name, err := c.resolveFile(pageFile)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	out := &pageSchema{}
	err = yaml.Unmarshal(data, out)
	if err != nil {
		return nil, err
	}

	actionLen := c.dev.Columns * c.dev.Rows

	pageObj := &page{
		path:    name,
		Timeout: out.Timeout,
		Actions: make([]actions.Action, actionLen),
	}

	for _, actionSchema := range out.Actions {
		actionObj, err := actions.LoadAction(actionSchema.ActionName, &actionSchema.Parameters)
		if err != nil {
			return nil, err
		}
		pageObj.Actions[actionSchema.Button] = actionObj
	}

	return pageObj, nil
}

func (c *controller) SwapPage(pageFile string) error {
	pageObj, err := c.resolvePage(pageFile)
	if err != nil {
		return err
	}

	c.pageWait.Lock()
	defer c.pageWait.Unlock()

	c.pageStack[len(c.pageStack)-1] = pageObj
	c.pageTop = pageObj

	return nil
}

func (c *controller) PushPage(pageFile string) error {
	pageObj, err := c.resolvePage(pageFile)
	if err != nil {
		return err
	}

	c.pageWait.Lock()
	defer c.pageWait.Unlock()

	c.pageStack = append(c.pageStack, pageObj)
	c.pageTop = pageObj

	return nil
}

func (c *controller) PopPage() error {
	c.pageWait.Lock()
	defer c.pageWait.Unlock()

	newLen := len(c.pageStack) - 1
	if newLen < 1 {
		return errors.New("page stack is empty")
	}

	c.pageStack = c.pageStack[:newLen]
	c.pageTop = c.pageStack[newLen-1]

	return nil
}
