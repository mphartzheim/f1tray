package processes

import (
	"crypto/sha256"
	"encoding/hex"
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

// hashCache stores computed SHA256 hashes for endpoints to detect changes.
var hashCache = make(map[string]string)

// LoadSchedule fetches and updates the schedule data and returns true if new data was loaded.
func LoadSchedule(url string, parseFunc func([]byte) (string, [][]string, error), status *widget.Label, tableContainer *fyne.Container) bool {
	body, title, rows, ok := fetchAndParse(url, parseFunc, status)
	if !ok {
		// No new data loaded.
		return false
	}

	highlightRow := -1
	var schedule models.ScheduleResponse
	err := json.Unmarshal(body, &schedule)
	if err != nil {
		status.SetText(fmt.Sprintf("Error parsing schedule: %v", err))
		return false
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
	return true
}

// LoadResults fetches and updates the race results and returns true if new data was loaded.
func LoadResults(url string, parseFunc func([]byte) (string, [][]string, error), status *widget.Label, tableContainer *fyne.Container) bool {
	_, title, rows, ok := fetchAndParse(url, parseFunc, status)
	if !ok {
		return false
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
	return true
}

// LoadUpcoming fetches and updates the upcoming race data and returns true if new data was loaded.
func LoadUpcoming(url string, parseFunc func([]byte) (string, [][]string, error), status *widget.Label, tableContainer *fyne.Container) bool {
	_, title, rows, ok := fetchAndParse(url, parseFunc, status)
	if !ok {
		return false
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
	return true
}

// fetchAndParse fetches JSON from a URL, compares its SHA256 hash to detect changes, and returns parsed data if updated.
func fetchAndParse(url string, parseFunc func([]byte) (string, [][]string, error), status *widget.Label) ([]byte, string, [][]string, bool) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		status.SetText(fmt.Sprintf("Fetch error: %v", err))
		return nil, "", nil, false
	}

	// Perform the request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		status.SetText(fmt.Sprintf("Fetch error: %v", err))
		return nil, "", nil, false
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		status.SetText(fmt.Sprintf("Read error: %v", err))
		return nil, "", nil, false
	}

	// Compute SHA256 hash of the body.
	newHashBytes := sha256.Sum256(body)
	newHash := hex.EncodeToString(newHashBytes[:])

	// Compare with the previously stored hash.
	if oldHash, ok := hashCache[url]; ok && oldHash == newHash {
		status.SetText("No new updates")
		return nil, "", nil, false
	}

	// Update the hash cache.
	hashCache[url] = newHash

	// Process the fetched data.
	title, rows, err := parseFunc(body)
	if err != nil {
		status.SetText(err.Error())
		return nil, "", nil, false
	}

	return body, title, rows, true
}
