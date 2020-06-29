package home

import (
	"sort"
	"strconv"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/common"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/parser_metadata"
	"github.com/spacetimi/pfh_reader_server/app_src/templates/colours"
	"github.com/spacetimi/pfh_reader_server/app_src/templates/graph_templates"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

type HomePageTab_t int

const (
	HOMEPAGE_TAB_DASHBOARD HomePageTab_t = iota
	HOMEPAGE_TAB_WEEK
	HOMEPAGE_TAB_INSIGHTS
	HOMEPAGE_TAB_SETTINGS
)

const HOMEPAGE_TABNAME_DASHBOARD = "dashboard"
const HOMEPAGE_TABNAME_WEEK = "week"
const HOMEPAGE_TABNAME_INSIGHTS = "insights"
const HOMEPAGE_TABNAME_SETTINGS = "settings"

func parseHomePageTabString(s string) HomePageTab_t {
	switch s {
	case HOMEPAGE_TABNAME_DASHBOARD:
		return HOMEPAGE_TAB_DASHBOARD
	case HOMEPAGE_TABNAME_WEEK:
		return HOMEPAGE_TAB_WEEK
	case HOMEPAGE_TABNAME_INSIGHTS:
		return HOMEPAGE_TAB_INSIGHTS
	case HOMEPAGE_TABNAME_SETTINGS:
		return HOMEPAGE_TAB_SETTINGS
	}
	return HOMEPAGE_TAB_DASHBOARD
}

func (tab HomePageTab_t) String() string {
	switch tab {
	case HOMEPAGE_TAB_DASHBOARD:
		return HOMEPAGE_TABNAME_DASHBOARD
	case HOMEPAGE_TAB_WEEK:
		return HOMEPAGE_TABNAME_WEEK
	case HOMEPAGE_TAB_INSIGHTS:
		return HOMEPAGE_TABNAME_INSIGHTS
	case HOMEPAGE_TAB_SETTINGS:
		return HOMEPAGE_TABNAME_SETTINGS
	}
	return HOMEPAGE_TABNAME_DASHBOARD
}

////////////////////////////////////////////////////////////////////////////////

const kPostArgNameCurrentTab = "tab"
const kPostArgNameCurrentDayIndex = "day-index"
const kPostArgNameCurrentWeekIndex = "week-index"
const kPostArgNameRuleIdToDelete = "rule-id-to-delete"

type parsedPostArgs struct {
	Tab              HomePageTab_t
	CurrentDayIndex  int // 0 is today, -1 is yesterday, and so on
	CurrentWeekIndex int // 0 is today, -1 is previous week, and so on
	RuleIdToDelete   int // -1 is unset
}

func parsePostArgs(postArgs map[string]string) *parsedPostArgs {
	parsed := &parsedPostArgs{
		Tab:              HOMEPAGE_TAB_DASHBOARD,
		CurrentDayIndex:  0,
		CurrentWeekIndex: 0,
		RuleIdToDelete:   -1,
	}

	if postArgs == nil || len(postArgs) == 0 {
		return parsed
	}

	currentTabName, ok := postArgs[kPostArgNameCurrentTab]
	if ok {
		parsed.Tab = parseHomePageTabString(currentTabName)
	}

	dayIndexString, ok := postArgs[kPostArgNameCurrentDayIndex]
	if ok {
		dayIndex, err := strconv.ParseInt(dayIndexString, 10, 32)
		if err != nil {
			logger.LogError("error parsing day index from post args" +
				"|day index string=" + dayIndexString +
				"|error=" + err.Error())
		} else {
			parsed.CurrentDayIndex = int(dayIndex)
		}
	}

	weekIndexString, ok := postArgs[kPostArgNameCurrentWeekIndex]
	if ok {
		weekIndex, err := strconv.ParseInt(weekIndexString, 10, 32)
		if err != nil {
			logger.LogError("error parsing week index from post args" +
				"|week index string=" + weekIndexString +
				"|error=" + err.Error())
		} else {
			parsed.CurrentWeekIndex = int(weekIndex)
		}
	}

	ruleIdToDeleteString, ok := postArgs[kPostArgNameRuleIdToDelete]
	if ok {
		ruleIdToDelete, err := strconv.ParseInt(ruleIdToDeleteString, 10, 32)
		if err != nil {
			logger.LogError("error parsing rule id to delete" +
				"|rule id to delete string=" + ruleIdToDeleteString +
				"|error=" + err.Error())
		} else {
			parsed.RuleIdToDelete = int(ruleIdToDelete)
		}
	}

	return parsed
}

func getColourForCategory(category app_core.Category_t) colours.Colour {
	switch category {
	case app_core.CATEGORY_PRODUCTIVE:
		return colours.MediumSeaGreen
	case app_core.CATEGORY_OPERATIONAL_OVERHEAD:
		return colours.DarkKhaki
	case app_core.CATEGORY_UNPRODUCTIVE:
		return colours.IndianRed
	case app_core.CATEGORY_UNCLASSIFIED:
		return colours.LightSteelBlue
	}

	return colours.LightSteelBlue
}

func formatTime(hours int, minutes int) string {
	hoursString := strconv.Itoa(hours % 12)
	minutesString := strconv.Itoa(minutes)
	if minutes <= 9 {
		minutesString = "0" + minutesString
	}
	suffix := "am"
	if hours > 12 {
		suffix = "pm"
	}

	return hoursString + ":" + minutesString + " " + suffix
}

func getHoursMinutesFromSeconds(seconds int) (int, int) {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60

	return hours, minutes
}

func getActivityOverviewAsBarGraph(activityOverviewData *common.ActivityOverviewData, graphName string) *graph_templates.BarGraphTemplateObject {

	datasets := make([]graph_templates.Dataset, app_core.CATEGORY_UNCLASSIFIED+1)

	for c := app_core.CATEGORY_PRODUCTIVE; c <= app_core.CATEGORY_UNCLASSIFIED; c = c + 1 {
		activity := activityOverviewData.GetActivityInPeriodsForCategory(c)
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
			GraphName:             graphName,
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

func getAppsUsageDatas(appUsageSecondsByAppName map[string]int64, maxAppsToShow int) []AppUsageData {
	appUsageDatas := make([]AppUsageData, 0)
	for appName, seconds := range appUsageSecondsByAppName {
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
	if len(appUsageDatas) > maxAppsToShow {
		appUsageDatas = appUsageDatas[0:maxAppsToShow]
	}

	return appUsageDatas
}
