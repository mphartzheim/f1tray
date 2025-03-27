package processes

import (
	"fmt"
	"strings"
	"time"
)

// Localize converts a UTC date and time string into local date and time strings.
func Localize(dateStr, timeStr string, use24h bool) (string, string) {
	combined := fmt.Sprintf("%sT%s", dateStr, timeStr)
	if !strings.HasSuffix(timeStr, "Z") &&
		!strings.Contains(timeStr, "+") &&
		!strings.Contains(timeStr, "-") {
		combined += "Z"
	}

	t, err := time.Parse(time.RFC3339, combined)
	if err != nil {
		fmt.Println("Error parsing datetime:", err)
		return dateStr, timeStr
	}

	local := t.Local()
	dateFormatted := local.Format("2006-01-02")
	timeFmt := "03:04 PM MST"
	if use24h {
		timeFmt = "15:04 MST"
	}
	return dateFormatted, local.Format(timeFmt)
}
