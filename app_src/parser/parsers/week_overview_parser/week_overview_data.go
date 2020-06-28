package week_overview_parser

import (
	"time"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/common"
)

type DayOfWeek int

const (
	Monday DayOfWeek = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
	NumDaysInWeek
)

type WeekOverviewData struct {
	common.WeekIdentifier

	WeekdaySummariesByDay map[DayOfWeek]*WeekdaySummaryData
}

type WeekdaySummaryData struct {
	LastUpdatedTime int64

	Day   int
	Month int
	Year  int
	DayOfWeek

	SecondsByCategory map[app_core.Category_t]int64
	SecondsByAppName  map[string]int64
	ActivityOverview  *common.ActivityOverviewData
}

////////////////////////////////////////////////////////////////////////////////

func (wod *WeekOverviewData) GetOrCreateSummaryForDay(year int, month int, day int) *WeekdaySummaryData {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	dow := GetDayOfWeekFromWeekday(t.Weekday())

	ws, ok := wod.WeekdaySummariesByDay[dow]
	if ok {
		return ws
	}

	ws = &WeekdaySummaryData{
		LastUpdatedTime: 0,

		Day:       day,
		Month:     month,
		Year:      year,
		DayOfWeek: dow,
	}
	wod.WeekdaySummariesByDay[dow] = ws

	return ws
}

func (wod *WeekOverviewData) GetTotalScreenTimeSeconds() int64 {
	totalSeconds := int64(0)

	for _, weekdaySummary := range wod.WeekdaySummariesByDay {
		for c := app_core.CATEGORY_PRODUCTIVE; c <= app_core.CATEGORY_UNCLASSIFIED; c = c + 1 {
			seconds, ok := weekdaySummary.SecondsByCategory[c]
			if ok {
				totalSeconds = totalSeconds + seconds
			}
		}
	}

	return totalSeconds
}

func (wod *WeekOverviewData) GetSecondsInCategory(category app_core.Category_t) int64 {
	totalSeconds := int64(0)

	for _, weekdaySummary := range wod.WeekdaySummariesByDay {
		seconds, ok := weekdaySummary.SecondsByCategory[category]
		if ok {
			totalSeconds = totalSeconds + seconds
		}
	}

	return totalSeconds
}

func (wod *WeekOverviewData) GetAppsUsageSeconds() map[string]int64 {
	result := make(map[string]int64)

	for _, weekdaySummary := range wod.WeekdaySummariesByDay {
		for appName, seconds := range weekdaySummary.SecondsByAppName {
			appUsage, ok := result[appName]
			if !ok {
				appUsage = 0
			}
			result[appName] = appUsage + seconds
		}
	}

	return result
}

func (wod *WeekOverviewData) GetSecondsInCategoryForDay(day DayOfWeek, category app_core.Category_t) int64 {
	wsd, ok := wod.WeekdaySummariesByDay[day]
	if !ok {
		return 0
	}

	seconds, ok := wsd.SecondsByCategory[category]
	if !ok {
		return 0
	}

	return seconds
}

func (wsd *WeekdaySummaryData) GetTotalScreentimeSeconds() int64 {
	totalSeconds := int64(0)

	for _, seconds := range wsd.SecondsByCategory {
		totalSeconds = totalSeconds + seconds
	}

	return totalSeconds
}

////////////////////////////////////////////////////////////////////////////////

func (dow DayOfWeek) String() string {
	switch dow {
	case Monday:
		return "Monday"
	case Tuesday:
		return "Tuesday"
	case Wednesday:
		return "Wednesday"
	case Thursday:
		return "Thursday"
	case Friday:
		return "Friday"
	case Saturday:
		return "Saturday"
	case Sunday:
		return "Sunday"
	default:
		return DayOfWeek(dow % NumDaysInWeek).String()
	}
}

/*
 time.Weekday starts from Sunday. We like to start from Monday
*/
func GetDayOfWeekFromWeekday(wd time.Weekday) DayOfWeek {
	corrected := (int(wd) + 6) % int(NumDaysInWeek)
	return DayOfWeek(corrected)
}
