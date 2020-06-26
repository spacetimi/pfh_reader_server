package common

import (
	"strconv"
	"time"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/timi_shared_server/utils/file_utils"
)

const kDAY_DATA_FILES_PREFIX = "day-"
const kDAY_DATA_FILES_EXTENSION = ".dat"

const kWEEK_DATA_FILES_PREFIX = "week-"
const kWEEK_DATA_FILES_EXTENSION = ".dat"

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

	fileName := kDAY_DATA_FILES_PREFIX + yearString + "-" + monthString + "-" + dayString + kDAY_DATA_FILES_EXTENSION
	return app_core.PFH_DAEMON_DATA_PATH + "/" + fileName
}

func GetRawDayDataFileNamesInDataFolder() ([]string, error) {
	filePaths, err := file_utils.GetFilePathsInDirectoryMatchingPattern(app_core.PFH_DAEMON_DATA_PATH,
		kDAY_DATA_FILES_PREFIX+"*"+kDAY_DATA_FILES_EXTENSION)
	if err != nil {
		return nil, err
	}

	return file_utils.GetFileNamesFromPaths(filePaths), nil
}

func GetWeekDataFilePath(weekIdentifier WeekIdentifier) string {
	fileName := kWEEK_DATA_FILES_PREFIX +
		strconv.Itoa(weekIdentifier.YearNumber) + "-" +
		strconv.Itoa(weekIdentifier.WeekNumber) +
		kWEEK_DATA_FILES_EXTENSION
	return app_core.PFH_DAEMON_DATA_PATH + "/" + fileName
}
