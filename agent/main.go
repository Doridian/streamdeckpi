package main

import (
	"bytes"
	"image/png"
	"log"

	"github.com/Doridian/streamdeck"

	_ "embed"
)

//go:embed test.png
var imageData []byte

func main() {
	img, err := png.Decode(bytes.NewReader(imageData))
	if err != nil {
		log.Panicf("No embdedded image: %v", err)
	}

	devs, err := streamdeck.Devices()
	if err != nil {
		log.Panicf("no Stream Deck devices found: %v", err)
	}
	if len(devs) == 0 {
		log.Panicf("no Stream Deck devices found")
	}
	d := devs[0]
	if err := d.Open(); err != nil {
		log.Panicf("can't open device: %v", err)
	}
	defer d.Close()

	d.SetImage(0, img)
}
