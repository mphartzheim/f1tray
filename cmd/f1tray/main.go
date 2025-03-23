package main

import (
	"fmt"
	"os"

	"f1tray/internal/tray"

	"github.com/getlantern/systray"
)

func main() {
	tray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("F1 Tray")
	systray.SetTooltip("F1 Session Notifier")

	// Load tray icon from file
	iconData, err := os.ReadFile("assets/tray_icon.png")
	if err == nil {
		systray.SetIcon(iconData)
	} else {
		fmt.Println("Failed to load icon:", err)
	}

	mQuit := systray.AddMenuItem("Quit", "Exit the application")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	fmt.Println("Exiting F1 Tray.")
	os.Exit(0)
}
