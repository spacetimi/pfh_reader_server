package home

import (
	"sort"
	"strconv"
	"time"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/common"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/week_overview_parser"
	"github.com/spacetimi/pfh_reader_server/app_src/templates/graph_templates"
	"github.com/spacetimi/timi_shared_server/utils/file_utils"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

const kMAX_TOP_APPS_TO_SHOW_PER_WEEK = 10

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
				ErrorString: "No data for " + currentWeekString,
			},
			CurrentWeekString:  currentWeekString,
			PrevWeekIndex:      prevWeekIndex,
			NextWeekIndex:      nextWeekIndex,
			ShowNextWeekButton: canShowNextWeekButton,
			ShowPrevWeekButton: canShowPrevWeekButton,
		}
		return weekviewPageObject
	}

	wop := &week_overview_parser.WeekOverviewParser{}
	wod, err := wop.ParseFile(dataFilePath)
	if err != nil {
		weekviewPageObject = &WeekviewData{
			ErrorablePage: ErrorablePage{
				HasError:    true,
				ErrorString: "Error parsing data for " + currentWeekString,
			},
			CurrentWeekString:  currentWeekString,
			PrevWeekIndex:      prevWeekIndex,
			NextWeekIndex:      nextWeekIndex,
			ShowNextWeekButton: canShowNextWeekButton,
			ShowPrevWeekButton: canShowPrevWeekButton,
		}
		return weekviewPageObject
	}

	averageActivity := wod.GetAverageActivityPeriods()
	totalScreentimeSeconds := wod.GetTotalScreenTimeSeconds()
	totalScreentimeHours, totalScreentimeMinutes := getHoursMinutesFromSeconds(int(totalScreentimeSeconds))

	appsUsage := wod.GetAppsUsageSeconds()
	appUsageDatas := getAppsUsageDatas(appsUsage, kMAX_TOP_APPS_TO_SHOW_PER_WEEK)

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

		TotalScreenTimeHours:   totalScreentimeHours,
		TotalScreenTimeMinutes: totalScreentimeMinutes,

		CategorySplitPieGraph:   *(getWeekCategorySplitAsPieGraph(wod)),
		AverageActivityBarGraph: *(getActivityOverviewAsBarGraph(averageActivity, "week-average-activity-bargraph")),
		WeekdayActivities:       getWeekdayActivities(wod),

		TopApps: appUsageDatas,
	}

	return weekviewPageObject
}

func getWeekCategorySplitAsPieGraph(wod *week_overview_parser.WeekOverviewData) *graph_templates.PieGraphTemplateObject {
	dataset := graph_templates.NewDataset()

	dataset.AddDataItem(float32(wod.GetSecondsInCategory(app_core.CATEGORY_PRODUCTIVE)), getColourForCategory(app_core.CATEGORY_PRODUCTIVE))
	dataset.AddDataItem(float32(wod.GetSecondsInCategory(app_core.CATEGORY_OPERATIONAL_OVERHEAD)), getColourForCategory(app_core.CATEGORY_OPERATIONAL_OVERHEAD))
	dataset.AddDataItem(float32(wod.GetSecondsInCategory(app_core.CATEGORY_UNPRODUCTIVE)), getColourForCategory(app_core.CATEGORY_UNPRODUCTIVE))
	dataset.AddDataItem(float32(wod.GetSecondsInCategory(app_core.CATEGORY_UNCLASSIFIED)), getColourForCategory(app_core.CATEGORY_UNCLASSIFIED))

	legends := []string{
		"Productive",
		"Operational Overhead",
		"Unproductive",
		"Others",
	}

	return &graph_templates.PieGraphTemplateObject{
		GraphTemplateObject: graph_templates.GraphTemplateObject{
			GraphName:             "week-category-split-piegraph",
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

func getWeekdayActivities(wod *week_overview_parser.WeekOverviewData) []WeekdayActivityData {
	var weekdayActivities []WeekdayActivityData

	for _, weekdaySummary := range wod.WeekdaySummariesByDay {
		weekdayActivity := WeekdayActivityData{
			WeekdayIndex:     int(weekdaySummary.DayOfWeek),
			WeekdayName:      weekdaySummary.DayOfWeek.String(),
			ActivityBarGraph: *(getActivityOverviewAsBarGraph(weekdaySummary.ActivityOverview, weekdaySummary.DayOfWeek.String()+"-activity-bargraph")),
		}
		weekdayActivities = append(weekdayActivities, weekdayActivity)
	}

	sort.Slice(weekdayActivities, func(i, j int) bool {
		return weekdayActivities[i].WeekdayIndex < weekdayActivities[j].WeekdayIndex
	})

	return weekdayActivities
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
