package home

import "github.com/spacetimi/pfh_reader_server/app_src/templates/graph_templates"

type HomePageObject struct {
	CurrentTab string

	DashboardData
	WeekviewData
}

type DashboardData struct {
	ErrorablePage

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

	TopApps []AppUsageData
}

type WeekviewData struct {
	ErrorablePage

	CurrentWeekString  string
	IsCurrentWeek      bool
	ShowPrevWeekButton bool
	ShowNextWeekButton bool
	PrevWeekIndex      int
	NextWeekIndex      int
}

type ErrorablePage struct {
	HasError    bool
	ErrorString string
}

type AppUsageData struct {
	AppName    string
	Seconds    int64
	TimeToShow string
}
