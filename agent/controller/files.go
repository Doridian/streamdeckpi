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

func (c *controller) cleanPath(file string) (string, error) {
	file = path.Clean(file)
	if file == "" || file == ".." || file[0] == '/' || (len(file) >= 3 && file[0:3] == "../") {
		return file, errors.New("paths outside config dir are not allowed")
	}
	return file, nil
}

func (c *controller) resolveFile(file string) (io.ReadCloser, error) {
	fh, err := os.Open(path.Join(loadMainDirectory, file))
	if err == nil {
		return fh, nil
	}

	if !os.IsNotExist(err) {
		return nil, err
	}

	reader, err := agent.FS.Open(file)
	return reader, err
}
