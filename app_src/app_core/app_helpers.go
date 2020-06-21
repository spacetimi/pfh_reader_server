package app_core

import (
	"strconv"
	"time"
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

	fileName := yearString + "-" + monthString + "-" + dayString + ".dat"
	return PFH_DAEMON_DATA_PATH + "/" + fileName
}
