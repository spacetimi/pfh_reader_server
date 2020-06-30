package collate

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/common"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/day_overview_parser"
	"github.com/spacetimi/pfh_reader_server/app_src/user_preferences"
	"github.com/spacetimi/timi_shared_server/utils/file_utils"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"github.com/spacetimi/timi_shared_server/utils/time_utils"
)

func CollateDaysToWeeks() error {

	if user_preferences.Instance() == nil {
		return errors.New("error getting user preferences")
	}

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

	deleteOldDayDataFiles(dayDataFileNames)

	return nil
}

func deleteOldDayDataFiles(fileNames []string) {
	timeNow := time.Now()

	for _, fileName := range fileNames {
		parsedFileName, err := parseDayDataFileName(fileName)
		if err != nil {
			logger.LogWarning("error parsing day-data file name" +
				"|file name=" + fileName +
				"|error=" + err.Error())
			continue
		}

		t := time.Date(parsedFileName.Year, time.Month(parsedFileName.Month), parsedFileName.Day, 0, 0, 0, 0, time.Local)
		age := time_utils.GetDurationBetweenTimes(t, timeNow)
		ageInDays := time_utils.DurationToDays(age)

		if ageInDays >= app_core.MAX_DAYS_TO_KEEP_RAW_DAY_DATA_FILES {
			err = os.Remove(app_core.PFH_DAEMON_DATA_PATH + "/" + fileName)
			if err != nil {
				logger.LogError("error removing stale day-data file" +
					"|file path=" + app_core.PFH_DAEMON_DATA_PATH + "/" + fileName +
					"|error=" + err.Error())
				continue
			}
		}
	}
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

	/*
	 If the day-data-file OR the user-preferences-data-file have been modified
	 after we collated this day into a week, reprocess it
	*/
	if dayDataFileModTime <= weekdaySummary.LastUpdatedTime &&
		user_preferences.Instance().DataModTime <= weekdaySummary.LastUpdatedTime {
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
