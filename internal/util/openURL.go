package util

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
)

func OpenWebPage(link string) {
	parsed, err := url.Parse(link)
	if err != nil {
		fmt.Println("❌ Invalid URL:", link)
		return
	}

	app := fyne.CurrentApp()
	if app == nil {
		fmt.Println("❌ Cannot open URL – app not initialized")
		return
	}

	err = app.OpenURL(parsed)
	if err != nil {
		fmt.Println("❌ Failed to open URL:", err)
	}
}
