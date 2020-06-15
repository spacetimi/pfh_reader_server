package home

import (
	"net/http"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/day_overview_parser"
	"github.com/spacetimi/pfh_reader_server/app_src/templates/graph_templates"
	"github.com/spacetimi/pfh_reader_server/app_src/templates/home_page_templates"

	"github.com/spacetimi/pfh_reader_server/app_src/app_routes"
	"github.com/spacetimi/timi_shared_server/code/config"
	"github.com/spacetimi/timi_shared_server/code/core/controller"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"github.com/spacetimi/timi_shared_server/utils/templated_writer"
)

type HomeHandler struct { // Implements IRouteHandler
	*templated_writer.TemplatedWriter
}

func NewHomeHandler() *HomeHandler {
	hh := &HomeHandler{}
	hh.TemplatedWriter = templated_writer.NewTemplatedWriter(config.GetAppTemplateFilesPath()+"/home", config.GetAppTemplateFilesPath()+"/graphs")

	// Parse templates for every request on LOCAL so that we can iterate over the templates
	// without having to restart the server every time
	hh.TemplatedWriter.ForceReparseTemplates = config.GetEnvironmentConfiguration().AppEnvironment == config.LOCAL

	return hh
}

func (hh *HomeHandler) Routes() []controller.Route {
	return []controller.Route{
		controller.NewRoute(app_routes.Home, []controller.RequestMethodType{controller.GET, controller.POST}),
		controller.NewRoute(app_routes.HomeSlash, []controller.RequestMethodType{controller.GET, controller.POST}),
	}
}

func (hh *HomeHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

	dop := &day_overview_parser.DayOverviewParser{}
	dod, e := dop.ParseFile("/Users/avkrishnan/Library/Containers/com.spacetimi.pfh-daemon/Data/Documents/2020-06-14.dat")
	if e != nil {
		logger.LogError(e.Error())
	} else {
		logger.VarDumpInfo("total time", dod.TotalTimeSeconds)
		logger.VarDumpInfo("productive time", dod.GetUsageSecondsInCategory(app_core.CATEGORY_PRODUCTIVE))
		logger.VarDumpInfo("unproductive time", dod.GetUsageSecondsInCategory(app_core.CATEGORY_UNPRODUCTIVE))
		logger.VarDumpInfo("operational-overhead time", dod.GetUsageSecondsInCategory(app_core.CATEGORY_OPERATIONAL_OVERHEAD))
		logger.VarDumpInfo("unclassified time", dod.GetUsageSecondsInCategory(app_core.CATEGORY_UNCLASSIFIED))
	}

	r := graph_templates.Color{R: 150, G: 80, B: 80, A: 0.9}
	g := graph_templates.Color{R: 80, G: 150, B: 80, A: 0.9}
	b := graph_templates.Color{R: 80, G: 80, B: 150, A: 0.9}

	pageObject := &home_page_templates.HomePageTemplate{
		PG: graph_templates.PieGraphTemplateObject{
			GraphTemplateObject: graph_templates.GraphTemplateObject{
				GraphName: "dummypiegraph",
				Dataset: graph_templates.Dataset{
					Data:    []float32{20, 20, 20},
					Colors:  []graph_templates.Color{},
					Legends: []string{"Productive", "Unproductive", "Operational Overhead"},
				},
				ShowLegend:     false,
				LegendPosition: "bottom",
			},
			IsDoughnut:       true,
			CutoutPercentage: 50,
		},
	}
	pageObject.PG.GraphTemplateObject.Dataset.Colors = append(pageObject.PG.GraphTemplateObject.Dataset.Colors, r)
	pageObject.PG.GraphTemplateObject.Dataset.Colors = append(pageObject.PG.GraphTemplateObject.Dataset.Colors, g)
	pageObject.PG.GraphTemplateObject.Dataset.Colors = append(pageObject.PG.GraphTemplateObject.Dataset.Colors, b)

	err := hh.TemplatedWriter.Render(httpResponseWriter, "home_page_template.html", pageObject)
	if err != nil {
		logger.LogError("error rendering home page" +
			"|error=" + err.Error())
		httpResponseWriter.WriteHeader(http.StatusInternalServerError)
	}
}
