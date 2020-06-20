package home_page_templates

import "github.com/spacetimi/pfh_reader_server/app_src/templates/graph_templates"

type HomePageTemplate struct {
	CategorySplitPieGraph graph_templates.PieGraphTemplateObject
	DailyActivityBarGraph graph_templates.BarGraphTemplateObject
}
