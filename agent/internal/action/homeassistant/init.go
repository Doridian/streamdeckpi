package homeassistant

import "github.com/Doridian/streamdeckpi/agent/internal/action/loader"

func init() {
	loader.RegisterAction(&haEntityAction{})
	loader.RegisterAction(&haLightAction{})
	loader.RegisterAction(&haStringAction{})
}
