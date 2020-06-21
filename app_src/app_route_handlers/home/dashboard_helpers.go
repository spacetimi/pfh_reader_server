package home

import (
	"net/http"
	"strconv"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/day_overview_parser"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/parser_metadata"
	"github.com/spacetimi/pfh_reader_server/app_src/templates/colours"
	"github.com/spacetimi/pfh_reader_server/app_src/templates/graph_templates"
	"github.com/spacetimi/pfh_reader_server/app_src/templates/home_page_templates"
	"github.com/spacetimi/timi_shared_server/code/core/controller"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

func (hh *HomeHandler) showDashboard(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs, postArgs *parsedPostArgs) {
	dop := &day_overview_parser.DayOverviewParser{}
	dod, e := dop.ParseFile(app_core.GetRawDayDataFilePath(postArgs.CurrentDayIndex))
	if e != nil {
		logger.LogError(e.Error())
		httpResponseWriter.WriteHeader(http.StatusInternalServerError)
	}

	pageObject := &home_page_templates.HomePageTemplate{
		CategorySplitPieGraph: *(getDayCategorySplitAsPieGraph(dod)),
		DailyActivityBarGraph: *(getDayActivityAsBarGraph(dod)),
	}

	err := hh.TemplatedWriter.Render(httpResponseWriter, "home_page_template.html", pageObject)
	if err != nil {
		logger.LogError("error rendering home page" +
			"|error=" + err.Error())
		httpResponseWriter.WriteHeader(http.StatusInternalServerError)
	}

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
			GraphName:         "today-category-split-piegraph",
			Datasets:          []graph_templates.Dataset{*dataset},
			Legends:           legends,
			ShowLegend:        true,
			LegendPosition:    "left",
			ResponsiveSize:    false,
			UseWidthAndHeight: false,
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
			dataset.AddDataItem(float32(int(seconds/60)), colour)
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
			GraphName:         "today-activity-bargraph",
			Datasets:          datasets,
			Legends:           legends,
			ShowLegend:        false,
			LegendPosition:    "top",
			UseWidthAndHeight: true,
			Width:             400,
			Height:            50,
			ResponsiveSize:    true,
		},
		Stacked:                   true,
		BarDisplayPercentage:      1.0,
		CategoryDisplayPercentage: 1.0,
		ShowAxis:                  false,
		ShowGridlines:             false,
		ShowTicks:                 false,
	}
}

// TODO: Move this somewhere else. Not just for dashboard
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

// TODO: Move this somewhere else. Not just for dashboard
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
