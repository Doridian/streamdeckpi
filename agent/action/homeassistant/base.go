package homeassistant

import (
	"fmt"
	"path"

	"github.com/Doridian/streamdeckpi/agent/controller"
	ha "github.com/kylegrantlucas/go-hass"
)

type haInstanceConfig struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

var haAccessors map[string]*ha.Access

func GetHomeAssistant(ctrl controller.Controller, id string) (*ha.Access, error) {
	if id == "" {
		id = "default"
	}

	access, ok := haAccessors[id]
	if !ok {
		config := &haInstanceConfig{}

		path := path.Join("/global/homeassistant", fmt.Sprintf("%s.yml", id))
		path, err := ctrl.CleanPath(path)
		if err != nil {
			return nil, err
		}

		err = ctrl.LoadConfig(path, config)
		if err != nil {
			return nil, err
		}

		access = ha.NewAccess(config.Url, config.Token)
	}

	return access, nil
}
