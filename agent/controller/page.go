package controller

import (
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
		actionObj, err := actions.LoadAction(actionSchema.ActionName, actionSchema.Parameters)
		if err != nil {
			return nil, err
		}
		pageObj.Actions[actionSchema.Button] = actionObj
	}

	return pageObj, nil
}

func (c *controller) SwapPage(pageFile string) error {
	return nil
}

func (c *controller) PushPage(pageFile string) error {
	return nil
}

func (c *controller) PopPage() error {
	return nil
}
