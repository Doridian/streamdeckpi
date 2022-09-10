package homeassistant

import "github.com/Doridian/streamdeckpi/agent/action/loader"

func init() {
	loader.RegisterAction(&haEntityAction{})
	loader.RegisterAction(&haLightAction{})
}
