package impl

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

var loadMainDirectory = loadMainDir()

func loadMainDir() string {
	configDir := os.Getenv("STREAMDECKPI_CONFIG_DIR")
	if configDir == "" {
		configDir = "/etc/streamdeckpi"
	}
	return path.Clean(configDir)
}

func (c *controllerImpl) CleanPath(file string) (string, error) {
	// This absolutizes paths and fixes definitons like "../../../"
	// which would refer outside the config dir automatically for us
	file = path.Clean(path.Join("/", file))
	if file[0] == '/' {
		// Strip out absolute paths after cleanup
		// this means they'll be relative to config dir
		// as opposed to realtive to current page dir or otherwise
		file = file[1:]
	}
	return file, nil
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

func (c *controllerImpl) LoadConfig(file string, v interface{}) error {
	reader, err := c.ResolveFile(file)
	if err != nil {
		return err
	}
	defer reader.Close()

	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)
	return decoder.Decode(v)
}

func (c *controllerImpl) ResolveFile(file string) (io.ReadCloser, error) {
	if file == "" {
		return nil, errors.New("refusing to open empty filename")
	}

	reader, err := os.Open(path.Join(loadMainDirectory, file))
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	return reader, assertFileIsFile(reader)
}
