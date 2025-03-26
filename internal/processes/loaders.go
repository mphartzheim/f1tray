package processes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"f1tray/internal/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// etagCache stores ETag values for endpoints to support conditional requests.
var etagCache = make(map[string]string)

func LoadSchedule(url string, parseFunc func([]byte) (string, [][]string, error), status *widget.Label, tableContainer *fyne.Container) {
	body, title, rows, ok := fetchAndParse(url, parseFunc, status)
	if !ok {
		// If no new data, exit early.
		return
	}

	highlightRow := -1
	var schedule models.ScheduleResponse
	err := json.Unmarshal(body, &schedule)
	if err != nil {
		status.SetText(fmt.Sprintf("Error parsing schedule: %v", err))
		return
	}

	now := time.Now()
	for i, race := range schedule.MRData.RaceTable.Races {
		raceDate, _ := time.Parse("2006-01-02", race.Date)
		if raceDate.After(now) || raceDate.Equal(now) {
			highlightRow = i + 1
			break
		}
	}

	table := widget.NewTable(
		func() (int, int) { return len(rows) + 1, 4 },
		func() fyne.CanvasObject {
			bg := canvas.NewRectangle(nil)
			label := widget.NewLabel("")
			return container.NewStack(bg, label)
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			wrapper := obj.(*fyne.Container)
			label := wrapper.Objects[1].(*widget.Label)
			bg := wrapper.Objects[0].(*canvas.Rectangle)
			if id.Row == 0 {
				headers := []string{"Round", "Race Name", "Circuit", "Location (Date)"}
				label.SetText(headers[id.Col])
				bg.Hide()
			} else {
				label.SetText(rows[id.Row-1][id.Col])
				if id.Row == highlightRow {
					bg.FillColor = theme.Color(theme.ColorNamePrimary)
					bg.Show()
				} else {
					bg.Hide()
				}
				bg.Resize(wrapper.Size())
			}
			wrapper.Refresh()
		},
	)

	table.SetColumnWidth(0, 60)
	table.SetColumnWidth(1, 200)
	table.SetColumnWidth(2, 280)
	table.SetColumnWidth(3, 280)
	table.Resize(fyne.NewSize(820, float32((len(rows)+1)*30)))

	tableContainer.Objects = []fyne.CanvasObject{table}
	tableContainer.Refresh()
	status.SetText(fmt.Sprintf("%s loaded", title))
}

func LoadResults(url string, parseFunc func([]byte) (string, [][]string, error), status *widget.Label, tableContainer *fyne.Container) {
	_, title, rows, ok := fetchAndParse(url, parseFunc, status)
	if !ok {
		return
	}

	table := widget.NewTable(
		func() (int, int) { return len(rows) + 1, 4 },
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			return container.New(layout.NewStackLayout(), label)
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*fyne.Container).Objects[0].(*widget.Label)
			if id.Row == 0 {
				headers := []string{"Pos", "Driver", "Team", "Time/Status"}
				label.SetText(headers[id.Col])
			} else {
				label.SetText(rows[id.Row-1][id.Col])
			}
		},
	)

	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 180)
	table.SetColumnWidth(2, 180)
	table.SetColumnWidth(3, 300)
	table.Resize(fyne.NewSize(600, float32((len(rows)+1)*30)))

	tableContainer.Objects = []fyne.CanvasObject{table}
	tableContainer.Refresh()
	status.SetText(fmt.Sprintf("Results loaded for %s", title))
}

func LoadUpcoming(url string, parseFunc func([]byte) (string, [][]string, error), status *widget.Label, tableContainer *fyne.Container) {
	_, title, rows, ok := fetchAndParse(url, parseFunc, status)
	if !ok {
		return
	}

	table := widget.NewTable(
		func() (int, int) { return len(rows) + 1, 3 },
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			return container.New(layout.NewStackLayout(), label)
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*fyne.Container).Objects[0].(*widget.Label)
			if id.Row == 0 {
				headers := []string{"Session", "Date", "Time"}
				label.SetText(headers[id.Col])
			} else {
				label.SetText(rows[id.Row-1][id.Col])
			}
		},
	)

	table.SetColumnWidth(0, 150)
	table.SetColumnWidth(1, 150)
	table.SetColumnWidth(2, 150)
	table.Resize(fyne.NewSize(500, float32((len(rows)+1)*30)))

	tableContainer.Objects = []fyne.CanvasObject{table}
	tableContainer.Refresh()
	status.SetText(fmt.Sprintf("Upcoming race loaded: %s", title))
}

func fetchAndParse(url string, parseFunc func([]byte) (string, [][]string, error), status *widget.Label) ([]byte, string, [][]string, bool) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		status.SetText(fmt.Sprintf("Fetch error: %v", err))
		return nil, "", nil, false
	}

	// Add the If-None-Match header if an ETag is cached for this URL.
	if etag, ok := etagCache[url]; ok && etag != "" {
		req.Header.Set("If-None-Match", etag)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		status.SetText(fmt.Sprintf("Fetch error: %v", err))
		return nil, "", nil, false
	}
	defer resp.Body.Close()

	// If the data hasn't changed, exit without updating.
	if resp.StatusCode == http.StatusNotModified {
		status.SetText("No new updates")
		return nil, "", nil, false
	}

	// Cache the new ETag if provided.
	newETag := resp.Header.Get("ETag")
	if newETag != "" {
		etagCache[url] = newETag
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		status.SetText(fmt.Sprintf("Read error: %v", err))
		return nil, "", nil, false
	}

	title, rows, err := parseFunc(body)
	if err != nil {
		status.SetText(err.Error())
		return nil, "", nil, false
	}

	return body, title, rows, true
}
