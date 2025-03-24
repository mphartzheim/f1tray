package main

import (
	"flag"
	"log"

	"f1tray/internal/config"
	"f1tray/internal/gui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

var debugMode bool

func main() {
	flag.BoolVar(&debugMode, "debug", false, "Enable debug mode")
	flag.Parse()

	f1App := app.NewWithID(config.AppID)

	gui.LoadAppIcon(f1App, config.IconPath)

	mainWindow := f1App.NewWindow("F1 Tray")
	mainWindow.Resize(fyne.NewSize(config.WindowWidth, config.WindowHeight))

	content := gui.BuildMainMenu(mainWindow, debugMode)
	mainWindow.SetContent(content)

	f1App.Settings().SetTheme(theme.DefaultTheme())

	mainWindow.SetCloseIntercept(func() {
		mainWindow.Hide()
	})

	f1App.Lifecycle().SetOnStopped(func() {
		log.Println("F1 Tray shutting down")
	})

	mainWindow.Show()
	mainWindow.SetMaster()
	f1App.Run()
}
