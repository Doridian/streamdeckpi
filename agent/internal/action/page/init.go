package page

import "github.com/Doridian/streamdeckpi/agent/internal/action/loader"

func init() {
	loader.RegisterAction(&popPage{})
	loader.RegisterAction(&pushPage{})
	loader.RegisterAction(&swapPage{})
}
