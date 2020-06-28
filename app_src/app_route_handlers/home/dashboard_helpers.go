package home

import (
	"sort"
	"strconv"
	"time"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/common"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/day_overview_parser"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/parser_metadata"
	"github.com/spacetimi/pfh_reader_server/app_src/templates/graph_templates"
	"github.com/spacetimi/timi_shared_server/utils/file_utils"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

const kMAX_TOP_APPS_TO_SHOW = 5

func (hh *HomeHandler) getDashboardPageObject(postArgs *parsedPostArgs) *DashboardData {

	var dashboardPageObject *DashboardData

	dataFilePath := common.GetRawDayDataFilePath(postArgs.CurrentDayIndex)

	if !file_utils.DoesFileOrDirectoryExist(dataFilePath) {
		dashboardPageObject = &DashboardData{
			CurrentDayString:  getCurrentDayStringFromDayIndex(postArgs.CurrentDayIndex),
			IsToday:           postArgs.CurrentDayIndex == 0,
			ShowPrevDayButton: -(postArgs.CurrentDayIndex) < app_core.MAX_DAYS_TO_KEEP_RAW_DAY_DATA_FILES,
			ShowNextDayButton: postArgs.CurrentDayIndex != 0,
			PrevDayIndex:      postArgs.CurrentDayIndex - 1,
			NextDayIndex:      postArgs.CurrentDayIndex + 1,

			ErrorablePage: ErrorablePage{
				HasError:    true,
				ErrorString: "No data for " + getCurrentDayStringFromDayIndex(postArgs.CurrentDayIndex),
			},
		}

		return dashboardPageObject
	}

	dop := &day_overview_parser.DayOverviewParser{}
	dod, e := dop.ParseFile(dataFilePath)
	if e != nil {
		dashboardPageObject = &DashboardData{
			CurrentDayString:  getCurrentDayStringFromDayIndex(postArgs.CurrentDayIndex),
			IsToday:           postArgs.CurrentDayIndex == 0,
			ShowPrevDayButton: -(postArgs.CurrentDayIndex) < app_core.MAX_DAYS_TO_KEEP_RAW_DAY_DATA_FILES,
			ShowNextDayButton: postArgs.CurrentDayIndex != 0,
			PrevDayIndex:      postArgs.CurrentDayIndex - 1,
			NextDayIndex:      postArgs.CurrentDayIndex + 1,

			ErrorablePage: ErrorablePage{
				HasError:    true,
				ErrorString: "Error parsing data for " + getCurrentDayStringFromDayIndex(postArgs.CurrentDayIndex),
			},
		}

		return dashboardPageObject
	}

	totalHours, totalMinutes := getHoursMinutesFromSeconds(int(dod.TotalTimeSeconds))

	appsUsage := dod.GetAppsUsageSeconds()
	appUsageDatas := make([]AppUsageData, 0)
	for appName, seconds := range appsUsage {
		hours, minutes := getHoursMinutesFromSeconds(int(seconds))
		timeToShow := ""
		if hours > 0 {
			timeToShow = strconv.Itoa(hours) + " hours "
		}
		if minutes > 0 {
			timeToShow = timeToShow + strconv.Itoa(minutes) + " min"
		}
		appUsageDatas = append(appUsageDatas, AppUsageData{
			AppName:    appName,
			Seconds:    seconds,
			TimeToShow: timeToShow,
		})
	}
	sort.Slice(appUsageDatas, func(i, j int) bool {
		return appUsageDatas[i].Seconds > appUsageDatas[j].Seconds
	})
	if len(appUsageDatas) > kMAX_TOP_APPS_TO_SHOW {
		appUsageDatas = appUsageDatas[0:kMAX_TOP_APPS_TO_SHOW]
	}

	dashboardPageObject = &DashboardData{
		CurrentDayString:  getCurrentDayStringFromDayIndex(postArgs.CurrentDayIndex),
		IsToday:           postArgs.CurrentDayIndex == 0,
		ShowPrevDayButton: -(postArgs.CurrentDayIndex) < app_core.MAX_DAYS_TO_KEEP_RAW_DAY_DATA_FILES,
		ShowNextDayButton: postArgs.CurrentDayIndex != 0,
		PrevDayIndex:      postArgs.CurrentDayIndex - 1,
		NextDayIndex:      postArgs.CurrentDayIndex + 1,

		ErrorablePage: ErrorablePage{
			HasError:    false,
			ErrorString: "",
		},

		TotalScreenTimeHours:   totalHours,
		TotalScreenTimeMinutes: totalMinutes,

		CategorySplitPieGraph: *(getDayCategorySplitAsPieGraph(dod)),
		DailyActivityBarGraph: *(getDayActivityAsBarGraph(dod)),

		TopApps: appUsageDatas,
	}

	return dashboardPageObject
}

func getDayCategorySplitAsPieGraph(dod *day_overview_parser.DayOverviewData) *graph_templates.PieGraphTemplateObject {
	dataset := graph_templates.NewDataset()
	dataset.AddDataItem(float32(dod.GetUsageSecondsInCategory(app_core.CATEGORY_PRODUCTIVE)), getColourForCategory(app_core.CATEGORY_PRODUCTIVE))
	dataset.AddDataItem(float32(dod.GetUsageSecondsInCategory(app_core.CATEGORY_OPERATIONAL_OVERHEAD)), getColourForCategory(app_core.CATEGORY_OPERATIONAL_OVERHEAD))
	dataset.AddDataItem(float32(dod.GetUsageSecondsInCategory(app_core.CATEGORY_UNPRODUCTIVE)), getColourForCategory(app_core.CATEGORY_UNPRODUCTIVE))
	dataset.AddDataItem(float32(dod.GetUsageSecondsInCategory(app_core.CATEGORY_UNCLASSIFIED)), getColourForCategory(app_core.CATEGORY_UNCLASSIFIED))

	legends := []string{
		"Productive",
		"Operational Overhead",
		"Unproductive",
		"Others",
	}

	return &graph_templates.PieGraphTemplateObject{
		GraphTemplateObject: graph_templates.GraphTemplateObject{
			GraphName:             "today-category-split-piegraph",
			Datasets:              []graph_templates.Dataset{*dataset},
			Legends:               legends,
			ShowLegend:            true,
			LegendPosition:        "left",
			ResponsiveSize:        false,
			UseWidthAndHeight:     false,
			FormatTimeFromSeconds: true,
		},
		IsDoughnut:       true,
		CutoutPercentage: 50,
	}
}

func getDayActivityAsBarGraph(dod *day_overview_parser.DayOverviewData) *graph_templates.BarGraphTemplateObject {

	datasets := make([]graph_templates.Dataset, app_core.CATEGORY_UNCLASSIFIED+1)

	for c := app_core.CATEGORY_PRODUCTIVE; c <= app_core.CATEGORY_UNCLASSIFIED; c = c + 1 {
		activity := dod.GetActivityInPeriodsForCategory(c)
		dataset := graph_templates.NewDataset()
		colour := getColourForCategory(c)
		for _, seconds := range activity {
			dataset.AddDataItem(float32(seconds), colour)
		}
		datasets[c] = *dataset
	}

	legends := make([]string, parser_metadata.NUM_ACTIVITY_PERIODS_PER_DAY)
	for i := 0; i < parser_metadata.NUM_ACTIVITY_PERIODS_PER_DAY; i = i + 1 {
		hours, minutes := parser_metadata.ParseActivityPeriodIndex(i)
		legends[i] = formatTime(hours, minutes)
	}

	return &graph_templates.BarGraphTemplateObject{
		GraphTemplateObject: graph_templates.GraphTemplateObject{
			GraphName:             "today-activity-bargraph",
			Datasets:              datasets,
			Legends:               legends,
			ShowLegend:            false,
			LegendPosition:        "top",
			UseWidthAndHeight:     true,
			Width:                 400,
			Height:                50,
			ResponsiveSize:        true,
			FormatTimeFromSeconds: true,
		},
		Stacked:                   true,
		BarDisplayPercentage:      1.0,
		CategoryDisplayPercentage: 1.0,
		ShowAxis:                  false,
		ShowGridlines:             false,
		ShowTicks:                 false,
	}
}

func getCurrentDayStringFromDayIndex(dayIndex int) string {
	if dayIndex > 0 {
		logger.LogError("invalid day index|day index=" + strconv.Itoa(dayIndex))
		return "invalid"
	}

	if dayIndex == 0 {
		return "Today"
	}

	t := time.Now().AddDate(0, 0, dayIndex)
	return t.Month().String() + " " + strconv.Itoa(t.Day())
}
