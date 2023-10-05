package impl

import (
	"time"

	"gopkg.in/yaml.v3"
)

type actionSchema struct {
	ActionName string    `yaml:"name"`
	Button     [2]int    `yaml:"button"`
	Parameters yaml.Node `yaml:"parameters"`
}

type pageSchema struct {
	Timeout time.Duration  `yaml:"timeout"`
	Actions []actionSchema `yaml:"actions"`
}
