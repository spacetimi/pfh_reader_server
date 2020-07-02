package graph_templates

type BarGraphTemplateObject struct {
	GraphTemplateObject

	Stacked bool

	ShowAxis      bool
	ShowGridlines bool
	ShowTicks     bool

	BarDisplayPercentage      float32
	CategoryDisplayPercentage float32

	YAxisSuggestedMax int
}
