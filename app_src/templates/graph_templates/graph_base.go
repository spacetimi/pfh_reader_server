package graph_templates

import "github.com/spacetimi/pfh_reader_server/app_src/templates/colours"

type GraphTemplateObject struct {
	GraphName string
	Dataset
	ShowLegend     bool
	LegendPosition string
}

type Dataset struct {
	Data    []float32
	Colours []colours.Colour
	Legends []string
}

func NewDataset() *Dataset {
	return &Dataset{}
}

func (d *Dataset) AddDataItem(dataItem float32, colour colours.Colour, legend string) {
	d.Data = append(d.Data, dataItem)
	d.Colours = append(d.Colours, colour)
	d.Legends = append(d.Legends, legend)
}
