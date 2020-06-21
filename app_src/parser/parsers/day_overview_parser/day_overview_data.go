package day_overview_parser

import (
	"strconv"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/parser_metadata"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"github.com/spacetimi/timi_shared_server/utils/time_utils"
)

type DayOverviewData struct {
	CategoryOverviewsByCategory map[app_core.Category_t]*CategoryOverviewData
	TotalTimeSeconds            int64
	ActivityOverview            *ActivityOverviewData
}

type CategoryOverviewData struct {
	Category               app_core.Category_t
	AppUsageOverviewsByApp map[string]*AppUsageOverviewData
	TotalTimeSeconds       int64
}

type AppUsageOverviewData struct {
	AppName          string
	TotalTimeSeconds int64
}

type ActivityOverviewData struct {
	ActivityPeriods [parser_metadata.NUM_ACTIVITY_PERIODS_PER_DAY]*ActivityPeriodData
}

type ActivityPeriodData struct {
	SecondsInCategory map[app_core.Category_t]int64
}

////////////////////////////////////////////////////////////////////////////////

func NewDayOverviewData() *DayOverviewData {
	return &DayOverviewData{
		CategoryOverviewsByCategory: make(map[app_core.Category_t]*CategoryOverviewData),
		TotalTimeSeconds:            0,
		ActivityOverview:            newActivityOverviewData(),
	}
}

func (dod *DayOverviewData) GetUsageSecondsInCategory(category app_core.Category_t) int64 {
	cod, ok := dod.CategoryOverviewsByCategory[category]
	if !ok {
		return 0
	}

	return cod.TotalTimeSeconds
}

func (dod *DayOverviewData) AddAppUsageSecondsInCategory(category app_core.Category_t, appName string, timestamp int64, seconds int64) {

	cod, ok := dod.CategoryOverviewsByCategory[category]
	if !ok {
		cod = newCategoryOverviewData(category)
		dod.CategoryOverviewsByCategory[category] = cod
	}
	cod.addAppUsageInSeconds(appName, seconds)

	dod.ActivityOverview.addActivity(category, timestamp, seconds)

	dod.TotalTimeSeconds += seconds
}

func (dod *DayOverviewData) GetActivityInPeriodsForCategory(category app_core.Category_t) []int64 {
	return dod.ActivityOverview.getActivityInPeriodsForCategory(category)
}

func (dod *DayOverviewData) GetAppsUsageSeconds() map[string]int64 {
	result := make(map[string]int64)

	for _, cod := range dod.CategoryOverviewsByCategory {
		for appName, appUsage := range cod.AppUsageOverviewsByApp {

			seconds, ok := result[appName]
			if !ok {
				result[appName] = 0
				seconds = 0
			}
			result[appName] = seconds + appUsage.TotalTimeSeconds
		}
	}

	return result
}

////////////////////////////////////////////////////////////////////////////////

func (cod *CategoryOverviewData) addAppUsageInSeconds(appName string, seconds int64) {
	appUsageOverview, ok := cod.AppUsageOverviewsByApp[appName]
	if !ok {
		appUsageOverview = newAppUsageOverviewData(appName)
		cod.AppUsageOverviewsByApp[appName] = appUsageOverview
	}
	appUsageOverview.addUsageInSeconds(seconds)

	cod.TotalTimeSeconds += seconds
}

func newCategoryOverviewData(category app_core.Category_t) *CategoryOverviewData {
	return &CategoryOverviewData{
		Category:               category,
		AppUsageOverviewsByApp: make(map[string]*AppUsageOverviewData),
		TotalTimeSeconds:       0,
	}
}

////////////////////////////////////////////////////////////////////////////////

func (appud *AppUsageOverviewData) addUsageInSeconds(seconds int64) {
	appud.TotalTimeSeconds += seconds
}

func newAppUsageOverviewData(appName string) *AppUsageOverviewData {
	return &AppUsageOverviewData{
		AppName:          appName,
		TotalTimeSeconds: 0,
	}
}

////////////////////////////////////////////////////////////////////////////////

func newActivityOverviewData() *ActivityOverviewData {

	aod := &ActivityOverviewData{}

	for i := 0; i < len(aod.ActivityPeriods); i = i + 1 {
		aod.ActivityPeriods[i] = newActivityPeriodData()
	}

	return aod
}

func newActivityPeriodData() *ActivityPeriodData {
	apd := &ActivityPeriodData{
		SecondsInCategory: make(map[app_core.Category_t]int64),
	}

	for i := app_core.CATEGORY_PRODUCTIVE; i <= app_core.CATEGORY_UNCLASSIFIED; i = i + 1 {
		apd.SecondsInCategory[i] = 0
	}

	return apd
}

func (aod *ActivityOverviewData) addActivity(category app_core.Category_t, timestampOfActivity int64, seconds int64) {
	timeOfActivity := time_utils.UnixTimeStampToLocalTime(timestampOfActivity)
	beginningOfDay := time_utils.BeginningOfDay(timeOfActivity)

	secondsIntoDay := time_utils.DurationToSeconds(time_utils.GetDurationBetweenTimes(beginningOfDay, timeOfActivity))
	periodIndex := int(secondsIntoDay / parser_metadata.ACTIVITY_PERIOD_LENGTH_SECONDS)

	if periodIndex < 0 || periodIndex >= len(aod.ActivityPeriods) {
		logger.LogError("period index out of range" +
			"|period index=" + strconv.Itoa(periodIndex) +
			"|timestamp of activity=" + strconv.FormatInt(timestampOfActivity, 10))
		return
	}
	activityPeriod := aod.ActivityPeriods[periodIndex]
	activityPeriod.addActivity(category, seconds)
}

func (apd *ActivityPeriodData) addActivity(category app_core.Category_t, seconds int64) {
	secondsInCategory, ok := apd.SecondsInCategory[category]
	if !ok {
		secondsInCategory = 0
		apd.SecondsInCategory[category] = secondsInCategory
	}

	apd.SecondsInCategory[category] = secondsInCategory + seconds
}

func (apd *ActivityPeriodData) getSecondsInCategory(category app_core.Category_t) int64 {
	secondsInCategory, ok := apd.SecondsInCategory[category]
	if !ok {
		return 0
	}
	return secondsInCategory
}

func (aod *ActivityOverviewData) getActivityInPeriodsForCategory(category app_core.Category_t) []int64 {
	result := make([]int64, len(aod.ActivityPeriods))

	for i, apd := range aod.ActivityPeriods {
		result[i] = apd.getSecondsInCategory(category)
	}

	return result
}
