package home

import (
	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/day_overview_parser"
	"github.com/spacetimi/pfh_reader_server/app_src/templates/colours"
	"github.com/spacetimi/pfh_reader_server/app_src/templates/graph_templates"
)

func getDayCategorySplitAsPieGraph(dod *day_overview_parser.DayOverviewData) *graph_templates.PieGraphTemplateObject {
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

func getDayActivityAsBarGraph() *graph_templates.BarGraphTemplateObject {

	datasetProductive := graph_templates.NewDataset()
	datasetProductive.AddDataItem(10, colours.MediumSeaGreen)
	datasetProductive.AddDataItem(0, colours.MediumSeaGreen)
	datasetProductive.AddDataItem(20, colours.MediumSeaGreen)

	datasetUnproductive := graph_templates.NewDataset()
	datasetUnproductive.AddDataItem(0, colours.IndianRed)
	datasetUnproductive.AddDataItem(20, colours.IndianRed)
	datasetUnproductive.AddDataItem(0, colours.IndianRed)

	legends := []string{
		"10 - 10:30",
		"10.30 - 11",
		"11 - 11:30",
	}

	return &graph_templates.BarGraphTemplateObject{
		GraphTemplateObject: graph_templates.GraphTemplateObject{
			GraphName:         "today-activity-bargraph",
			Datasets:          []graph_templates.Dataset{*datasetProductive, *datasetUnproductive},
			Legends:           legends,
			ShowLegend:        false,
			LegendPosition:    "top",
			UseWidthAndHeight: true,
			Width:             400,
			Height:            40,
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
