package about

import (
    "github.com/spacetimi/pfh_reader_server/app_src/app_routes"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "github.com/spacetimi/timi_shared_server/utils/templated_writer"
    "net/http"
)

type AboutHandler struct {        // Implements IRouteHandler
    *templated_writer.TemplatedWriter
}

func NewAboutHandler() *AboutHandler {
    ah := &AboutHandler{}
    ah.TemplatedWriter = templated_writer.NewTemplatedWriter(config.GetAppTemplateFilesPath() + "/about")

    // Parse templates for every request on LOCAL so that we can iterate over the templates
    // without having to restart the server every time
    ah.TemplatedWriter.ForceReparseTemplates = config.GetEnvironmentConfiguration().AppEnvironment == config.LOCAL

    return ah
}

func (ah *AboutHandler) Routes() []controller.Route {
    return []controller.Route{
        controller.NewRoute(app_routes.About, []controller.RequestMethodType{controller.GET, controller.POST}),
    }
}

func (ah *AboutHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

    err := ah.TemplatedWriter.Render(httpResponseWriter, "about_page_template.html", nil)
    if err != nil {
        logger.LogError("error rendering about page" +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
    }
}
