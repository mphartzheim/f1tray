package util

import "fyne.io/fyne/v2/container"

func SelectTabByName(tabs *container.AppTabs, title string) {
	for i, item := range tabs.Items {
		if item.Text == title {
			tabs.SelectIndex(i)
			return
		}
	}
}
