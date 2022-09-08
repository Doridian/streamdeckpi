package controller

import (
	"errors"
	"io"
	"os"
	"path"

	"github.com/Doridian/streamdeckpi/agent"
)

var loadMainDirectory = loadMainDir()

func loadMainDir() string {
	configDir := os.Getenv("STREAMDECKPI_CONFIG_DIR")
	if configDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		configDir = path.Join(cwd, "config")
	}
	return configDir
}

func (c *controller) resolveFile(file string) (io.ReadCloser, error) {
	if path.IsAbs(file) {
		return nil, errors.New("absolute paths are not allowed")
	}

	fh, err := os.Open(path.Join(loadMainDirectory, file))
	if err == nil {
		return fh, nil
	}

	if !os.IsNotExist(err) {
		return nil, err
	}

	return agent.FS.Open(file)
}
