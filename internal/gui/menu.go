package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func BuildMainMenu(mainWindow fyne.Window, debugMode bool) fyne.CanvasObject {
	label := widget.NewLabel("Welcome to F1 Tray!")

	content := container.NewVBox(
		label,
		SessionButton("View Full Race Schedule", "schedule", mainWindow),
		widget.NewSeparator(),
		SessionButton("Race Results", "results", mainWindow),
		SessionButton("Qualifying Results", "qualifying", mainWindow),
		SessionButton("Sprint Qualifying Results", "sprint-qualifying", mainWindow),
		SessionButton("Sprint Results", "sprint", mainWindow),
		SessionButton("Practice 1 Results", "practice/1", mainWindow),
		SessionButton("Practice 2 Results", "practice/2", mainWindow),
		SessionButton("Practice 3 Results", "practice/3", mainWindow),
		widget.NewSeparator(),
	)

	if debugMode {
		content.Add(widget.NewLabel("[DEBUG MODE ENABLED]"))
		content.Add(widget.NewButton("Run Debug Test", func() {
			fmt.Println("Debug button clicked")
		}))
	}

	content.Add(widget.NewButton("Quit", func() {
		fyne.CurrentApp().Quit()
	}))

	return container.NewVScroll(content)
}

func SessionButton(label, session string, win fyne.Window) *widget.Button {
	if session == "schedule" {
		return widget.NewButton(label, func() {
			view, err := BuildRaceScheduleView(win, BuildMainMenu(win, false))
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			win.SetContent(view)
		})
	}
	return widget.NewButton(label, func() {
		view, err := BuildSessionResults(win, session, BuildMainMenu(win, false))
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		win.SetContent(view)
	})
}
