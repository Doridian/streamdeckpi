package misc

import (
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/Doridian/streamdeck"
	"github.com/Doridian/streamdeckpi/agent/actions"
	"github.com/Doridian/streamdeckpi/agent/controller"
)

type emptyStruct struct{}

type Command struct {
	actions.ActionBase

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

func (a *Command) setCurrentIcon(icon string) {
	if icon == "" {
		icon = a.Icon
	}
	if icon == a.currentIcon {
		return
	}

	a.currentIcon = icon
	a.doRender = true
}

func (a *Command) ApplyConfig(imageLoader controller.ImageLoader, controller controller.Controller) error {
	err := a.ActionBase.ApplyConfig(imageLoader, controller)
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

func (a *Command) Run(pressed bool) error {
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

func (a *Command) Name() string {
	return "command"
}

func (a *Command) Render(force bool) (*streamdeck.ImageData, error) {
	if force || a.doRender {
		a.doRender = false
		return a.ImageLoader.Load(a.currentIcon)
	}
	return nil, nil
}
