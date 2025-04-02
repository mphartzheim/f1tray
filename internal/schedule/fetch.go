package schedule

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mphartzheim/f1tray/internal/appstate"
	"github.com/mphartzheim/f1tray/internal/models"
)

func FetchSchedule(state *appstate.AppState) ([]ScheduledRace, error) {
	resp, err := http.Get(fmt.Sprintf(models.ScheduleURL, state.SelectedYear))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res ScheduleAPIResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return res.MRData.RaceTable.Races, nil
}
