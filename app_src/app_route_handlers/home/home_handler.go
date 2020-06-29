package home

import (
	"net/http"

	"github.com/spacetimi/pfh_reader_server/app_src/app_routes"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/collate"
	"github.com/spacetimi/pfh_reader_server/app_src/user_preferences"
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

	postArgs := parsePostArgs(args.PostArgs)

	if postArgs.RuleIdToDelete >= 0 {
		hh.deleteRule(postArgs.RuleIdToDelete)
		user_preferences.Reload()
	}

	if postArgs.AddRuleArgs != nil {
		hh.addRule(postArgs.AddRuleArgs.MatchType, postArgs.AddRuleArgs.MatchExpression,
			postArgs.AddRuleArgs.ShouldMatchCase, postArgs.AddRuleArgs.Category)
		user_preferences.Reload()
	}

	/* Try to collate data on every request so we can keep up */
	collate.CollateDaysToWeeks()

	dashboardPageObject := hh.getDashboardPageObject(postArgs)
	weekviewPageObject := hh.getWeekviewPageObject(postArgs)
	settingsPageObject := hh.getSettingsPageObject()

	pageObject := HomePageObject{
		CurrentTab: postArgs.Tab.String(),

		DashboardData: *dashboardPageObject,
		WeekviewData:  *weekviewPageObject,
		SettingsData:  *settingsPageObject,
	}

	err := hh.TemplatedWriter.Render(httpResponseWriter, "home_page_template.html", pageObject)
	if err != nil {
		logger.LogError("error rendering home page" +
			"|error=" + err.Error())
		httpResponseWriter.WriteHeader(http.StatusInternalServerError)
	}
}
