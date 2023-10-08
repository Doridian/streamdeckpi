package misc

import (
	"log"
	"sync"
	"time"

	"github.com/Doridian/go-streamdeck"
	"github.com/Doridian/streamdeckpi/agent/internal/action"
	"github.com/Doridian/streamdeckpi/agent/internal/action/loader"
	"github.com/Doridian/streamdeckpi/agent/internal/controller"
	"gopkg.in/yaml.v3"
)

type multiActionCondition struct {
	Pressed     bool          `yaml:"pressed"`
	MinDuration time.Duration `yaml:"min"`
	MaxDuration time.Duration `yaml:"max"`
}

type multiActionConfig struct {
	Conditions []multiActionCondition `yaml:"conditions"`
	Pressed    bool                   `yaml:"pressed"`
	Name       string                 `yaml:"name"`
	Parameters yaml.Node              `yaml:"parameters"`
}

type multiActionMapped struct {
	conditions []multiActionCondition
	pressed    bool
	action     action.Action
}

type runEntry struct {
	time    time.Time
	pressed bool
}

type multi struct {
	action.ActionBase

	RunActionConfig    []*multiActionConfig `yaml:"run"`
	RenderActionConfig *mapActionConfig     `yaml:"render"`

	runActions            []*multiActionMapped
	longestRunActionChain int

	renderAction action.Action

	runHistory       []*runEntry
	lastRunCheck     time.Time
	lastPressedState bool

	checkActionLock sync.Mutex
}

func (m *multiActionMapped) matches(runHistory []*runEntry) bool {
	if len(runHistory) < len(m.conditions) {
		return false
	}

	runHistoryOffset := len(runHistory) - len(m.conditions)

	var nextActionTime time.Time
	for i, condition := range m.conditions {
		historyEntry := runHistory[runHistoryOffset+i]
		if historyEntry.pressed != condition.Pressed {
			return false
		}

		if i < len(m.conditions)-1 {
			nextActionTime = runHistory[runHistoryOffset+i+1].time
		} else {
			nextActionTime = time.Now()
		}

		actionDuration := nextActionTime.Sub(historyEntry.time)

		if actionDuration < condition.MinDuration {
			return false
		}

		if condition.MaxDuration > 0 && actionDuration > condition.MaxDuration {
			return false
		}
	}

	return true
}

func (a *multi) New() action.Action {
	return &multi{
		runHistory: make([]*runEntry, 0),
	}
}

func (a *multi) ApplyConfig(config *yaml.Node, imageHelper controller.ImageHelper, ctrl controller.Controller) error {
	err := a.ActionBase.ApplyConfig(config, imageHelper, ctrl)
	if err != nil {
		return err
	}

	err = config.Decode(a)
	if err != nil {
		return err
	}

	a.runActions = make([]*multiActionMapped, 0, len(a.RunActionConfig))
	for _, runActionConfig := range a.RunActionConfig {
		runAction, err := loader.LoadAction(runActionConfig.Name, &runActionConfig.Parameters, imageHelper, ctrl)
		if err != nil {
			return err
		}
		a.runActions = append(a.runActions, &multiActionMapped{
			action:     runAction,
			pressed:    runActionConfig.Pressed,
			conditions: runActionConfig.Conditions,
		})

		runActionChain := len(runActionConfig.Conditions)
		if runActionChain > a.longestRunActionChain {
			a.longestRunActionChain = runActionChain
		}
	}

	a.renderAction, err = loader.LoadAction(a.RenderActionConfig.Name, &a.RenderActionConfig.Parameters, imageHelper, ctrl)
	if err != nil {
		return err
	}

	return nil
}

func (a *multi) checkActions() error {
	if !a.checkActionLock.TryLock() {
		return nil
	}
	defer a.checkActionLock.Unlock()

	a.lastRunCheck = time.Now()

	for _, runAction := range a.runActions {
		if runAction.matches(a.runHistory) {
			a.runHistory = make([]*runEntry, 0)
			return runAction.action.Run(runAction.pressed)
		}
	}

	return nil
}

func (a *multi) Run(pressed bool) error {
	if a.lastPressedState == pressed {
		return nil
	}
	a.lastPressedState = pressed

	a.runHistory = append(a.runHistory, &runEntry{
		time:    time.Now(),
		pressed: pressed,
	})
	if len(a.runHistory) > a.longestRunActionChain {
		a.runHistory = a.runHistory[1:]
	}

	return a.checkActions()
}

func (a *multi) goCheckActions() {
	err := a.checkActions()
	if err != nil {
		log.Printf("Error checking actions: %v", err)
	}
}

func (a *multi) Render(force bool) (*streamdeck.ImageData, error) {
	if time.Since(a.lastRunCheck) > time.Millisecond*10 {
		a.lastRunCheck = time.Now()
		go a.goCheckActions()
	}
	return a.renderAction.Render(force)
}

func (a *multi) Name() string {
	return "multi"
}
