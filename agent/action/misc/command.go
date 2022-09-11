package misc

import (
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/Doridian/go-streamdeck"
	"github.com/Doridian/streamdeckpi/agent/action"
	"github.com/Doridian/streamdeckpi/agent/controller"
	"gopkg.in/yaml.v3"
)

type emptyStruct struct{}

type command struct {
	action.ActionBase

	Command string   `yaml:"command"`
	Args    []string `yaml:"args"`

	Icon            string         `yaml:"icon"`
	RunningIcon     string         `yaml:"running_icon"`
	ExitCodeIcons   map[int]string `yaml:"exit_code_icons"`
	ExitDefaultIcon string         `yaml:"exit_default_icon"`
	ExitToIdleTime  time.Duration  `yaml:"exit_to_idle_time"`

	currentIcon string
	doRender    bool

	running   bool
	runSymbol *emptyStruct
	runLock   sync.Mutex
}

func (a *command) New() action.Action {
	return &command{}
}

func (a *command) setCurrentIcon(icon string) {
	if icon == "" {
		icon = a.Icon
	}
	if icon == a.currentIcon {
		return
	}

	a.currentIcon = icon
	a.doRender = true
}

func (a *command) ApplyConfig(config *yaml.Node, imageHelper controller.ImageHelper, ctrl controller.Controller) error {
	err := a.ActionBase.ApplyConfig(config, imageHelper, ctrl)
	if err != nil {
		return err
	}

	err = a.ApplyConfig(config, imageHelper, ctrl)
	if err != nil {
		return err
	}

	if a.ExitCodeIcons == nil {
		a.ExitCodeIcons = make(map[int]string)
	}

	a.currentIcon = a.Icon
	a.doRender = true
	return nil
}

func (a *command) Run(pressed bool) error {
	if !pressed {
		return nil
	}

	localRunSymbol := &emptyStruct{}

	a.runLock.Lock()
	if a.running {
		a.runLock.Unlock()
		return nil
	}
	a.running = true
	a.runSymbol = localRunSymbol
	a.runLock.Unlock()

	cmd := exec.Command(a.Command, a.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	a.setCurrentIcon(a.RunningIcon)

	cmdErr := cmd.Run()
	exitCode := 0
	if cmdErr != nil {
		exitError, ok := cmdErr.(*exec.ExitError)
		if ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = -1
		}
	}

	unlockRun := func() {
		a.runLock.Lock()
		a.running = false
		a.runLock.Unlock()
	}

	icon, hadExitCodeHandler := a.ExitCodeIcons[exitCode]
	if !hadExitCodeHandler {
		icon = a.ExitDefaultIcon
	}
	a.setCurrentIcon(icon)
	unlockRun()

	go func() {
		time.Sleep(a.ExitToIdleTime)

		a.runLock.Lock()
		if a.runSymbol == localRunSymbol {
			a.setCurrentIcon(a.Icon)
		}
		a.runLock.Unlock()
	}()

	if hadExitCodeHandler {
		return nil
	}
	return cmdErr
}

func (a *command) Name() string {
	return "command"
}

func (a *command) Render(force bool) (*streamdeck.ImageData, error) {
	if force || a.doRender {
		a.doRender = false
		return a.ImageHelper.Load(a.currentIcon)
	}
	return nil, nil
}
