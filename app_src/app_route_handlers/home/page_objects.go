package home

import (
	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/templates/graph_templates"
)

type HomePageObject struct {
	CurrentTab string

	DashboardData
	WeekviewData
	SettingsData
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

	WeekdayActivitiesBarGraph graph_templates.BarGraphTemplateObject
	WeekdayActivities         []WeekdayActivityData

	TopApps []AppUsageData
}

type WeekdayActivityData struct {
	WeekdayIndex int

	ScreentimeHours   int
	ScreentimeMinutes int

	WeekdayName      string
	ActivityBarGraph graph_templates.BarGraphTemplateObject
}

type SettingsData struct {
	ErrorablePage

	AppNameMatchRules     []SettingsMatchRule
	AppTitleBarMatchRules []SettingsMatchRule
}

type SettingsMatchRule struct {
	RuleId int

	MatchExpression string
	Category        app_core.Category_t
	ShouldMatchCase bool
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
