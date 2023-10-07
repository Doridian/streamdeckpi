package impl

import (
	"errors"
	"log"

	"github.com/Doridian/streamdeckpi/agent/internal/action"
)

func (c *controllerImpl) buttonLoop() {
	defer c.runWait.Done()

	buttonChan, err := c.dev.ReadKeys()
	if err != nil {
		log.Panicf("can't read keys: %v", err)
	}

	for c.running {
		b, ok := <-buttonChan
		if !ok {
			break
		}

		go c.handleButtonPress(int(b.Index), b.Pressed)
	}

	c.stopError(errors.New("unexpectedly reached end of button loop"))
}

func (c *controllerImpl) handleButtonPress(idx int, pressed bool) {
	actionObj := c.pageTop.actions[idx]
	c.runActionHandleError(actionObj, pressed)
}

func (c *controllerImpl) runActionHandleError(actionObj action.Action, pressed bool) {
	if actionObj == nil {
		return
	}
	err := actionObj.Run(pressed)
	if err != nil {
		log.Printf("Error running action: %v", err)
	}
}
