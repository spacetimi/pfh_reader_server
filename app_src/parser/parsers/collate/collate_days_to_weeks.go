package collate

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/common"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/day_overview_parser"
	"github.com/spacetimi/timi_shared_server/utils/file_utils"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

func CollateDaysToWeeks() error {

	dayDataFileNames, err := common.GetRawDayDataFileNamesInDataFolder()
	if err != nil {
		return errors.New("error getting day data file names: " + err.Error())
	}

	for _, dayDatadayDataFileName := range dayDataFileNames {
		err = processDayDataFile(dayDatadayDataFileName)
		if err != nil {
			logger.LogError("error processing day data file" +
				"|file name=" + dayDatadayDataFileName +
				"|error=" + err.Error())
			continue
		}
	}

	WeekOverviewCacheInstance().Apply()

	return nil
}

func processDayDataFile(fileName string) error {
	parsedFileName, err := parseDayDataFileName(fileName)
	if err != nil {
		return errors.New("error parsing day-data file name: " + err.Error())
	}

	weekIdentifier := getWeekForDay(parsedFileName.Year, parsedFileName.Month, parsedFileName.Day)
	weekOverview, err := WeekOverviewCacheInstance().GetWeekOverview(weekIdentifier)
	if err != nil {
		return errors.New("error getting week overview: " + err.Error())
	}

	weekdaySummary := weekOverview.GetOrCreateSummaryForDay(parsedFileName.Year, parsedFileName.Month, parsedFileName.Day)

	dayDataFilePath := app_core.PFH_DAEMON_DATA_PATH + "/" + fileName
	dayDataFileModTime, err := file_utils.GetFileModTimeUnix(dayDataFilePath)
	if err != nil {
		return errors.New("error getting file mtime for day data file at: " + dayDataFilePath)
	}

	if dayDataFileModTime <= weekdaySummary.LastUpdatedTime {
		return nil
	}

	weekdaySummary.LastUpdatedTime = time.Now().Unix()

	dop := &day_overview_parser.DayOverviewParser{}
	dod, err := dop.ParseFile(dayDataFilePath)

	weekdaySummary.SecondsByAppName = dod.GetAppsUsageSeconds()

	weekdaySummary.SecondsByCategory = make(map[app_core.Category_t]int64)
	for cat := app_core.CATEGORY_PRODUCTIVE; cat <= app_core.CATEGORY_UNCLASSIFIED; cat = cat + 1 {
		weekdaySummary.SecondsByCategory[cat] = dod.GetUsageSecondsInCategory(cat)
	}

	weekdaySummary.ActivityOverview = dod.ActivityOverview

	return nil
}

func parseDayDataFileName(fileName string) (*parsedDayDataFileName, error) {
	tokens := strings.Split(file_utils.GetFileNameWithoutExtension(fileName), "-")
	if len(tokens) != 4 {
		return nil, errors.New("malformed day-data file name: " + fileName)
	}

	day, err_day := strconv.Atoi(tokens[3])
	month, err_month := strconv.Atoi(tokens[2])
	year, err_year := strconv.Atoi(tokens[1])
	if err_day != nil || err_month != nil || err_year != nil {
		return nil, errors.New("malformed day-data file name: " + fileName)
	}

	parsed := parsedDayDataFileName{
		Day:   day,
		Month: month,
		Year:  year,
	}
	return &parsed, nil
}

type parsedDayDataFileName struct {
	Day   int
	Month int
	Year  int
}

func getWeekForDay(year int, month int, day int) common.WeekIdentifier {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	yearNumber, weekNumber := t.ISOWeek()

	return common.WeekIdentifier{
		YearNumber: yearNumber,
		WeekNumber: weekNumber,
	}
}
