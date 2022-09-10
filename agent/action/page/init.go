package page

import "github.com/Doridian/streamdeckpi/agent/action/loader"

func init() {
	loader.RegisterAction(&PopPage{})
	loader.RegisterAction(&PushPage{})
	loader.RegisterAction(&SwapPage{})
}
