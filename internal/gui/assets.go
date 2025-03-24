package gui

import (
	"io"
	"log"
	"os"

	"fyne.io/fyne/v2"
)

// LoadAppIcon loads the application icon from a file and applies it to the given Fyne app.
func LoadAppIcon(app fyne.App, path string) {
	iconFile, err := os.Open(path)
	if err != nil {
		log.Printf("Failed to load icon: %v", err)
		return
	}
	defer iconFile.Close()

	iconData, err := io.ReadAll(iconFile)
	if err != nil {
		log.Printf("Failed to read icon: %v", err)
		return
	}

	resource := fyne.NewStaticResource("tray_icon.png", iconData)
	app.SetIcon(resource)
}
