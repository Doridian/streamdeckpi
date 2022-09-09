package impl

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/Doridian/streamdeckpi/agent"
)

var loadMainDirectory = loadMainDir()

func loadMainDir() string {
	configDir := os.Getenv("STREAMDECKPI_CONFIG_DIR")
	if configDir == "" {
		configDir = "config"
	}
	return path.Clean(configDir)
}

func (c *controllerImpl) cleanPath(file string) (string, error) {
	return path.Clean(file), nil
}

func assertFileIsFile(file fs.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if stat.IsDir() {
		file.Close()
		return errors.New("tried opening a directory")
	}
	return nil
}

func (c *controllerImpl) resolveFile(file string) (io.ReadCloser, error) {
	if file == "" {
		return nil, errors.New("refusing to open empty filename")
	}

	var err error
	var reader fs.File

	reader, err = os.Open(path.Join(loadMainDirectory, file))
	if err == nil {
		return reader, assertFileIsFile(reader)
	}

	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("could not open real file: %w", err)
	}

	if file[0] == '/' { // Strip out absolute paths, they don't work in embed
		file = file[1:]
	}

	reader, err = agent.FS.Open(path.Join("embed", file))
	if err != nil {
		return nil, fmt.Errorf("could not open embedded file: %w", err)
	}
	return reader, assertFileIsFile(reader)
}
