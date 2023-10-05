package controller

import "io"

type Controller interface {
	SwapPage(pageFile string) error
	PushPage(pageFile string) error
	PopPage() error

	Start() error
	Stop() error
	Wait() error
	Reset() error

	SetBrightness(brightness int) error

	ResolveFile(file string) (io.ReadCloser, error)
	CleanPath(file string) (string, error)
	LoadConfig(file string, v interface{}) error
}
