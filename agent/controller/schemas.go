package controller

import (
	"time"
)

type actionSchema struct {
	ActionName string                 `yaml:"name"`
	Button     int                    `yaml:"button"`
	Parameters map[string]interface{} `yaml:"parameters"`
}

type pageSchema struct {
	path string

	Timeout time.Duration  `yaml:"timeout"`
	Actions []actionSchema `yaml:"actions"`
}
