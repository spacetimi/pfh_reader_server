package common

import (
	"strconv"
	"time"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
)

/**
day-index of 0 is today
day-index of -1 is yesterday, and so on
day-index with positive values is undefined
*/
func GetRawDayDataFilePath(dayIndex int) string {
	t := time.Now().AddDate(0, 0, dayIndex)

	year := t.Year()
	month := int(t.Month())
	day := t.Day()

	yearString := strconv.Itoa(year)
	monthString := strconv.Itoa(month)
	dayString := strconv.Itoa(day)

	if month <= 9 {
		monthString = "0" + monthString
	}
	if day <= 9 {
		dayString = "0" + dayString
	}

	fileName := "day-" + yearString + "-" + monthString + "-" + dayString + ".dat"
	return app_core.PFH_DAEMON_DATA_PATH + "/" + fileName
}
