package main

import (
	"log"

	"github.com/Doridian/go-streamdeck"
	controller_impl "github.com/Doridian/streamdeckpi/agent/controller/impl"

	_ "embed"
)

func main() {
	devs, err := streamdeck.Devices()
	if err != nil {
		log.Panicf("no Stream Deck devices found: %v", err)
	}
	if len(devs) == 0 {
		log.Panicf("no Stream Deck devices found")
	}

	controller, err := controller_impl.NewController(&devs[0])
	if err != nil {
		log.Panicf("Error initializing: %v", err)
	}

	err = controller.Start()
	if err != nil {
		log.Panicf("Error starting: %v", err)
	}
	defer controller.Stop()

	err = controller.Wait()
	if err != nil {
		log.Panicf("Error running: %v", err)
	}
}
