package day_overview_parser

import (
	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/common"
)

type DayOverviewData struct {
	CategoryOverviewsByCategory map[app_core.Category_t]*CategoryOverviewData
	TotalTimeSeconds            int64
	ActivityOverview            *common.ActivityOverviewData
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

////////////////////////////////////////////////////////////////////////////////

func NewDayOverviewData() *DayOverviewData {
	return &DayOverviewData{
		CategoryOverviewsByCategory: make(map[app_core.Category_t]*CategoryOverviewData),
		TotalTimeSeconds:            0,
		ActivityOverview:            common.NewActivityOverviewData(),
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

	dod.ActivityOverview.AddActivity(category, timestamp, seconds)

	dod.TotalTimeSeconds += seconds
}

func (dod *DayOverviewData) GetActivityInPeriodsForCategory(category app_core.Category_t) []int64 {
	return dod.ActivityOverview.GetActivityInPeriodsForCategory(category)
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
