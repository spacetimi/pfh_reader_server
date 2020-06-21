package app_core

import (
	"strconv"

	"github.com/spacetimi/timi_shared_server/utils/time_utils"
)

/**
day-index of 0 is today
day-index of -1 is yesterday, and so on
day-index with positive values is undefined
*/
func GetRawDayDataFilePath(dayIndex int) string {
	year := time_utils.GetLocalYear()
	month := time_utils.GetLocalMonth()
	day := time_utils.GetLocalDayOfMonth()

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
