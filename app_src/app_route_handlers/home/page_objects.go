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
	ShowPrevWeekButton bool
	ShowNextWeekButton bool
	PrevWeekIndex      int
	NextWeekIndex      int

	TotalScreenTimeHours   int
	TotalScreenTimeMinutes int

	CategorySplitPieGraph graph_templates.PieGraphTemplateObject

	AverageActivityBarGraph graph_templates.BarGraphTemplateObject
	WeekdayActivities       []WeekdayActivityData

	TopApps []AppUsageData
}

type WeekdayActivityData struct {
	WeekdayIndex int

	ScreentimeHours   int
	ScreentimeMinutes int

	WeekdayName      string
	ActivityBarGraph graph_templates.BarGraphTemplateObject
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
