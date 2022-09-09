package impl

import (
	"errors"
	"log"
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
	action := c.pageTop.actions[idx]
	if action == nil {
		return
	}
	err := action.Run(pressed)
	if err != nil {
		log.Printf("Error running action: %v", err)
	}
}
