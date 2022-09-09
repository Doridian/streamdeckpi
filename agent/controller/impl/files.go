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
	var err error
	var reader fs.File

	reader, err = os.Open(path.Join(loadMainDirectory, file))
	if err == nil {
		return reader, assertFileIsFile(reader)
	}

	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("could not open real file: %w", err)
	}

	reader, err = agent.FS.Open(path.Join("embed", file))
	if err != nil {
		return nil, fmt.Errorf("could not open embedded file: %w", err)
	}
	return reader, assertFileIsFile(reader)
}
