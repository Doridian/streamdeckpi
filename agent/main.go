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

//go:embed test2.png
var image2Data []byte

func main() {
	img, err := png.Decode(bytes.NewReader(imageData))
	if err != nil {
		log.Panicf("No embdedded image: %v", err)
	}

	img2, err := png.Decode(bytes.NewReader(image2Data))
	if err != nil {
		log.Panicf("No embdedded image2: %v", err)
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
	d.SetImage(1, img2)

	buttonChan, err := d.ReadKeys()
	if err != nil {
		log.Panicf("can't read keys: %v", err)
	}

	for {
		b := <-buttonChan
		log.Printf("B: %v, P: %v", b.Index, b.Pressed)
		if b.Pressed {
			d.SetImage(b.Index, img)
		} else {
			d.SetImage(b.Index, img2)
		}
	}
}
