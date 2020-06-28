package home

import (
	"time"

	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/common"
	"github.com/spacetimi/timi_shared_server/utils/file_utils"
)

func (hh *HomeHandler) getWeekviewPageObject(postArgs *parsedPostArgs) *WeekviewData {

	var weekviewPageObject *WeekviewData

	weekIdentifier := getWeekIdentifierFromWeekIndex(postArgs.CurrentWeekIndex)
	dataFilePath := common.GetWeekDataFilePath(weekIdentifier)

	if !file_utils.DoesFileOrDirectoryExist(dataFilePath) {
		weekviewPageObject = &WeekviewData{
			ErrorablePage: ErrorablePage{
				HasError:    true,
				ErrorString: "No data for week #placeholder#",
			},
			CurrentWeekString: "placeholder",
		}
		return weekviewPageObject
	}

	return &WeekviewData{}
}

/*
 week-index of 0 is current week, -1 is previous week, and so on
 positive values are undefined
*/
func getWeekIdentifierFromWeekIndex(weekIndex int) common.WeekIdentifier {
	t := time.Now()
	offsetWeeks := -1 * time.Duration(weekIndex) * 3600 * 24 * 7 * time.Second

	t.Add(offsetWeeks)

	yearNumber, weekNumber := t.ISOWeek()
	return common.WeekIdentifier{
		WeekNumber: weekNumber,
		YearNumber: yearNumber,
	}
}
