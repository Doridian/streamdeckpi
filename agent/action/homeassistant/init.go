package homeassistant

import "github.com/Doridian/streamdeckpi/agent/action/loader"

func init() {
	loader.RegisterAction(&HAEntityAction{})
}