package misc

import "github.com/Doridian/streamdeckpi/agent/action/loader"

func init() {
	loader.RegisterAction(&brightness{})
	loader.RegisterAction(&command{})
	loader.RegisterAction(&exit{})
	loader.RegisterAction(&reMap{})
	loader.RegisterAction(&multi{})
	loader.RegisterAction(&none{})
}
