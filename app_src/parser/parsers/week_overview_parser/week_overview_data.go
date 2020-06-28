package week_overview_parser

import (
	"time"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/common"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/parser_metadata"
	"github.com/spacetimi/timi_shared_server/utils/time_utils"
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

func (wod *WeekOverviewData) GetAverageActivityPeriods() *common.ActivityOverviewData {

	averageActivityOverview := common.NewActivityOverviewData()

	for _, weekdaySummary := range wod.WeekdaySummariesByDay {
		for _, activityPeriod := range weekdaySummary.ActivityOverview.ActivityPeriods {
			timestamp := time_utils.BeginningOfDay(time.Now()).Unix() + int64(activityPeriod.PeriodIndex*parser_metadata.ACTIVITY_PERIOD_LENGTH_SECONDS)

			averageActivityOverview.AddActivity(app_core.CATEGORY_PRODUCTIVE, timestamp, activityPeriod.GetSecondsInCategory(app_core.CATEGORY_PRODUCTIVE))
			averageActivityOverview.AddActivity(app_core.CATEGORY_OPERATIONAL_OVERHEAD, timestamp, activityPeriod.GetSecondsInCategory(app_core.CATEGORY_OPERATIONAL_OVERHEAD))
			averageActivityOverview.AddActivity(app_core.CATEGORY_UNPRODUCTIVE, timestamp, activityPeriod.GetSecondsInCategory(app_core.CATEGORY_UNPRODUCTIVE))
			averageActivityOverview.AddActivity(app_core.CATEGORY_UNCLASSIFIED, timestamp, activityPeriod.GetSecondsInCategory(app_core.CATEGORY_UNCLASSIFIED))
		}
	}

	numDaysWithData := int64(len(wod.WeekdaySummariesByDay))

	for _, activityPeriod := range averageActivityOverview.ActivityPeriods {
		for cat := app_core.CATEGORY_PRODUCTIVE; cat <= app_core.CATEGORY_UNCLASSIFIED; cat = cat + 1 {
			activityPeriod.SecondsInCategory[cat] = activityPeriod.GetSecondsInCategory(cat) / numDaysWithData
		}
	}

	return averageActivityOverview
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
