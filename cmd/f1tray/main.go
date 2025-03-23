package main

import (
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var debugMode bool

func main() {
	// Parse debug flag
	flag.BoolVar(&debugMode, "debug", false, "Enable debug mode")
	flag.Parse()

	// Create app
	f1App := app.NewWithID("com.f1tray.app")
	f1App.Settings().SetTheme(theme.LightTheme())

	// Load window icon
	iconFile, err := os.Open("assets/tray_icon.png")
	if err != nil {
		log.Printf("Failed to load icon: %v", err)
	} else {
		img, err := png.Decode(iconFile)
		if err != nil {
			log.Printf("Failed to decode icon: %v", err)
		} else {
			resource := fyne.NewStaticResource("icon.png", nil)
			f1App.SetIcon(fyne.NewStaticResource("icon.png", nil))
			f1App.SetIcon(resource)
			f1App.SetIcon(fyne.NewStaticResource("icon.png", nil))
			f1App.Settings().SetTheme(theme.LightTheme())
		}
		iconFile.Close()
	}

	// Create window
	mainWindow := f1App.NewWindow("F1 Tray")
	mainWindow.Resize(fyne.NewSize(400, 300))

	// Main content
	content := container.NewVBox(
		widget.NewLabel("Welcome to F1 Tray!"),
		widget.NewButton("Fetch Results (Coming Soon)", func() {
			fmt.Println("Fetch button pressed")
		}),
	)

	if debugMode {
		content.Add(widget.NewLabel("[DEBUG MODE ENABLED]"))
		content.Add(widget.NewButton("Run Debug Test", func() {
			fmt.Println("Debug button clicked")
		}))
	}

	mainWindow.SetContent(content)
	mainWindow.ShowAndRun()
}
