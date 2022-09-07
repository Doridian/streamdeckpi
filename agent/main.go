package main

import (
	"fmt"

	"github.com/muesli/streamdeck"
)

type programState struct {
	d streamdeck.Device
}

func (s *programState) closeStreamDeck() error {
	return s.d.Close()
}

func (s *programState) initStreamDeck() error {
	devs, err := streamdeck.Devices()
	if err != nil {
		return fmt.Errorf("no Stream Deck devices found: %s", err)
	}
	if len(devs) == 0 {
		return fmt.Errorf("no Stream Deck devices found")
	}
	s.d = devs[0]

	if err := s.d.Open(); err != nil {
		return fmt.Errorf("can't open device: %s", err)
	}

	/*
		ver, err := d.FirmwareVersion()
		if err != nil {
			return fmt.Errorf("can't retrieve device info: %s", err)
		}
		fmt.Printf("Found device with serial %s (firmware %s)\n",
			d.Serial, ver)
	*/

	return nil
}

func main() {

}
