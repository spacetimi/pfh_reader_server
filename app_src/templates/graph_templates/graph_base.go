package graph_templates

type GraphTemplateObject struct {
	GraphName string
	Dataset
	ShowLegend     bool
	LegendPosition string
}

type Dataset struct {
	Data    []float32
	Colors  []Color
	Legends []string
}

type Color struct {
	R int16
	G int16
	B int16
	A float32
}
