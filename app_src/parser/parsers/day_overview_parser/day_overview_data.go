package day_overview_parser

import (
	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
)

type DayOverviewData struct {
	CategoryOverviewsByCategory map[app_core.Category_t]*CategoryOverviewData
	TotalTimeSeconds            int64
}

func (dod *DayOverviewData) AddAppUsageSecondsInCategory(category app_core.Category_t, appName string, seconds int64) {
	cod, ok := dod.CategoryOverviewsByCategory[category]
	if !ok {
		cod = newCategoryOverviewData(category)
		dod.CategoryOverviewsByCategory[category] = cod
	}
	cod.addAppUsageInSeconds(appName, seconds)

	dod.TotalTimeSeconds += seconds
}

func NewDayOverviewData() *DayOverviewData {
	return &DayOverviewData{
		CategoryOverviewsByCategory: make(map[app_core.Category_t]*CategoryOverviewData),
		TotalTimeSeconds:            0,
	}
}

////////////////////////////////////////////////////////////////////////////////

type CategoryOverviewData struct {
	Category               app_core.Category_t
	AppUsageOverviewsByApp map[string]*AppUsageOverviewData
	TotalTimeSeconds       int64
}

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

type AppUsageOverviewData struct {
	AppName          string
	TotalTimeSeconds int64
}

func (appud *AppUsageOverviewData) addUsageInSeconds(seconds int64) {
	appud.TotalTimeSeconds += seconds
}

func newAppUsageOverviewData(appName string) *AppUsageOverviewData {
	return &AppUsageOverviewData{
		AppName:          appName,
		TotalTimeSeconds: 0,
	}
}
