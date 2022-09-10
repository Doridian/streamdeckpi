package misc

import "github.com/Doridian/streamdeckpi/agent/action/loader"

func init() {
	loader.RegisterAction(&Brightness{})
	loader.RegisterAction(&Command{})
	loader.RegisterAction(&Exit{})
	loader.RegisterAction(&Map{})
	loader.RegisterAction(&None{})
}
