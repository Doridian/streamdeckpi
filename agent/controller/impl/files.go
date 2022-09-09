package impl

import (
	"errors"
	"fmt"
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

func (c *controllerImpl) cleanPath(file string) (string, error) {
	file = path.Clean(file)
	if file == "" || file == ".." || file[0] == '/' || (len(file) >= 3 && file[0:3] == "../") {
		return file, errors.New("paths outside config dir are not allowed")
	}
	return file, nil
}

func (c *controllerImpl) resolveFile(file string) (io.ReadCloser, error) {
	fh, err := os.Open(path.Join(loadMainDirectory, file))
	if err == nil {
		return fh, nil
	}

	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("could not open real file: %w", err)
	}

	reader, err := agent.FS.Open(path.Join("embed", file))
	if err != nil {
		return nil, fmt.Errorf("could not open embedded file: %w", err)
	}
	return reader, nil
}
