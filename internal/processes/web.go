package processes

import (
	"net/url"

	"fyne.io/fyne/v2"
)

// OpenWebPage opens the specified URL in the default browser.
func OpenWebPage(link string) error {
	u, err := url.Parse(link)
	if err != nil {
		return err
	}
	return fyne.CurrentApp().OpenURL(u)
}
