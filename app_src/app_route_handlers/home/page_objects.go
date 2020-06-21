package home

import "github.com/spacetimi/pfh_reader_server/app_src/templates/graph_templates"

type HomePageObject struct {
	DashboardData
}

type DashboardData struct {
	HasError    bool
	ErrorString string

	CurrentDayString  string
	IsToday           bool
	ShowPrevDayButton bool
	ShowNextDayButton bool
	PrevDayIndex      int
	NextDayIndex      int

	TotalScreenTimeHours   int
	TotalScreenTimeMinutes int

	CategorySplitPieGraph graph_templates.PieGraphTemplateObject
	DailyActivityBarGraph graph_templates.BarGraphTemplateObject
}
