package impl

import (
	"time"

	"github.com/Doridian/streamdeckpi/agent/utils"
)

type actionSchema struct {
	ActionName string               `yaml:"name"`
	Button     [2]int               `yaml:"button"`
	Parameters utils.YAMLRawMessage `yaml:"parameters"`
}

type pageSchema struct {
	path string

	Timeout time.Duration  `yaml:"timeout"`
	Actions []actionSchema `yaml:"actions"`
}
