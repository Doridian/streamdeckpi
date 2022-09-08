package controller

import (
	"errors"
	"io"
)

func (c *controller) resolveFile(file string) (io.ReadCloser, error) {
	return nil, errors.New("could not find file")
}
