package home

import (
	"strconv"
	"time"

	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/common"
	"github.com/spacetimi/timi_shared_server/utils/file_utils"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

func (hh *HomeHandler) getWeekviewPageObject(postArgs *parsedPostArgs) *WeekviewData {

	var weekviewPageObject *WeekviewData

	weekIdentifier := getWeekIdentifierFromWeekIndex(postArgs.CurrentWeekIndex)
	dataFilePath := common.GetWeekDataFilePath(weekIdentifier)

	currentWeekString := getCurrentWeekStringFromWeekIndex(postArgs.CurrentWeekIndex)
	prevWeekIndex := postArgs.CurrentWeekIndex - 1
	nextWeekIndex := postArgs.CurrentWeekIndex + 1
	canShowNextWeekButton := postArgs.CurrentWeekIndex < 0
	canShowPrevWeekButton := true // TODO: Add clamp on max prev-weeks data to save?

	if !file_utils.DoesFileOrDirectoryExist(dataFilePath) {
		weekviewPageObject = &WeekviewData{
			ErrorablePage: ErrorablePage{
				HasError:    true,
				ErrorString: "No data for week #placeholder#",
			},
			CurrentWeekString:  currentWeekString,
			PrevWeekIndex:      prevWeekIndex,
			NextWeekIndex:      nextWeekIndex,
			ShowNextWeekButton: canShowNextWeekButton,
			ShowPrevWeekButton: canShowPrevWeekButton,
		}
		return weekviewPageObject
	}

	weekviewPageObject = &WeekviewData{
		ErrorablePage: ErrorablePage{
			HasError:    false,
			ErrorString: "",
		},

		CurrentWeekString:  currentWeekString,
		PrevWeekIndex:      prevWeekIndex,
		NextWeekIndex:      nextWeekIndex,
		ShowNextWeekButton: canShowNextWeekButton,
		ShowPrevWeekButton: canShowPrevWeekButton,
	}

	return weekviewPageObject
}

/*
 week-index of 0 is current week, -1 is previous week, and so on
 positive values are undefined
*/
func getWeekIdentifierFromWeekIndex(weekIndex int) common.WeekIdentifier {
	t := time.Now()
	t = t.AddDate(0, 0, 7*weekIndex)

	yearNumber, weekNumber := t.ISOWeek()
	return common.WeekIdentifier{
		WeekNumber: weekNumber,
		YearNumber: yearNumber,
	}
}

func getCurrentWeekStringFromWeekIndex(weekIndex int) string {
	if weekIndex > 0 {
		logger.LogError("invalid week index" + "|week index=" + strconv.Itoa(weekIndex))
		return "invalid"
	}

	if weekIndex == 0 {
		return "This Week"
	}

	t := time.Now().AddDate(0, 0, 7*weekIndex)
	t = getClosestPreviousMonday(t)

	return "Week of " + t.Month().String() + " " + strconv.Itoa(t.Day())
}

func getClosestPreviousMonday(t time.Time) time.Time {
	daysSinceMonday := (int(t.Weekday()-time.Monday) + 7) % 7
	t = t.AddDate(0, 0, -1*daysSinceMonday)
	return t
}
