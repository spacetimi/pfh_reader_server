package home

import (
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/day_overview_parser"
	"github.com/spacetimi/pfh_reader_server/app_src/templates/graph_templates"
	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/templates/colours"
)

func getPieGraphForDashboard(dod *day_overview_parser.DayOverviewData) *graph_templates.PieGraphTemplateObject {
	dataset := graph_templates.NewDataset()
	dataset.AddDataItem(float32(dod.GetUsageSecondsInCategory(app_core.CATEGORY_PRODUCTIVE)), colours.MediumSeaGreen)
	dataset.AddDataItem(float32(dod.GetUsageSecondsInCategory(app_core.CATEGORY_OPERATIONAL_OVERHEAD)), colours.DarkKhaki)
	dataset.AddDataItem(float32(dod.GetUsageSecondsInCategory(app_core.CATEGORY_UNPRODUCTIVE)), colours.IndianRed)
	dataset.AddDataItem(float32(dod.GetUsageSecondsInCategory(app_core.CATEGORY_UNCLASSIFIED)), colours.LightSteelBlue)

	legends := []string{
		"Productive",
		"Operational Overhead",
		"Unproductive",
		"Others",
	}

	return &graph_templates.PieGraphTemplateObject{
		GraphTemplateObject: graph_templates.GraphTemplateObject{
			GraphName:      "today-category-split-piegraph",
			Datasets:       []graph_templates.Dataset{*dataset},
			Legends:        legends,
			ShowLegend:     true,
			LegendPosition: "left",
		},
		IsDoughnut:       true,
		CutoutPercentage: 50,
	}
}