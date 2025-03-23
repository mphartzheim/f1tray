package tray

import (
	"github.com/getlantern/systray"
)

func Run(onReady func(), onExit func()) {
	systray.Run(onReady, onExit)
}
