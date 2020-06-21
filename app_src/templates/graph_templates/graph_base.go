package graph_templates

import "github.com/spacetimi/pfh_reader_server/app_src/templates/colours"

type GraphTemplateObject struct {
	GraphName string
	Datasets  []Dataset

	Legends        []string
	ShowLegend     bool
	LegendPosition string

	UseWidthAndHeight bool
	Width             int
	Height            int
	ResponsiveSize    bool

	FormatTimeFromSeconds bool
}

type Dataset struct {
	Data    []float32
	Colours []colours.Colour
}

func NewDataset() *Dataset {
	return &Dataset{}
}

func (d *Dataset) AddDataItem(dataItem float32, colour colours.Colour) {
	d.Data = append(d.Data, dataItem)
	d.Colours = append(d.Colours, colour)
}
