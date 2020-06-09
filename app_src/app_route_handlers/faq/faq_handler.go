package faq

import (
    "github.com/spacetimi/pfh_reader_server/app_src/app_routes"
    "github.com/spacetimi/pfh_reader_server/app_src/app_utils/app_simple_message_page"
    "github.com/spacetimi/pfh_reader_server/app_src/metadata/faq"
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service"
    "github.com/spacetimi/timi_shared_server/code/core/services/metadata_service/metadata_typedefs"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "github.com/spacetimi/timi_shared_server/utils/templated_writer"
    "net/http"
    "strconv"
)

type FaqHandler struct {        // Implements IRouteHandler
    *templated_writer.TemplatedWriter
}

func NewFaqHandler() *FaqHandler {
    fh := &FaqHandler{}
    fh.TemplatedWriter = templated_writer.NewTemplatedWriter(config.GetAppTemplateFilesPath() + "/faq")

    // Parse templates for every request on LOCAL so that we can iterate over the templates
    // without having to restart the server every time
    fh.TemplatedWriter.ForceReparseTemplates = config.GetEnvironmentConfiguration().AppEnvironment == config.LOCAL

    return fh
}

func (fh *FaqHandler) Routes() []controller.Route {
    return []controller.Route{
        controller.NewRoute(app_routes.Faq, []controller.RequestMethodType{controller.GET, controller.POST}),
    }
}

func (fh *FaqHandler) HandlerFunc(httpResponseWriter http.ResponseWriter, request *http.Request, args *controller.HandlerFuncArgs) {

    latestVersion, err := metadata_service.Instance().GetLatestDefinedVersion(metadata_typedefs.METADATA_SPACE_APP)
    if err != nil {
        logger.LogError("error getting latest metadata version" +
                        "|error=" + err.Error())
        // Show error message and return
        messageHeader := "Something went wrong"
        messageBody := "Please try again"
        backlinkName := "<< Home"
        app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
        return
    }

    faqMetadata := &faq.Metadata{}
    err = metadata_service.Instance().GetMetadataItem(faqMetadata, latestVersion)
    if err != nil {
        logger.LogError("error getting faq metadata" +
                        "|version=" + latestVersion.String() +
                        "|error=" + err.Error())
        // Show error message and return
        messageHeader := "Something went wrong"
        messageBody := "Please try again"
        backlinkName := "<< Home"
        app_simple_message_page.ShowAppSimpleMessagePage(httpResponseWriter, messageHeader, messageBody, app_routes.HomeSlash, backlinkName)
        return
    }

    pageObject := &PageObject{}
    for i, faqItem := range faqMetadata.FaqItems {
        pageObject.FaqItems = append(pageObject.FaqItems, Item{
            Id: "faq_item_" + strconv.Itoa(i),
            Question:faqItem.Question,
            Answer:faqItem.Answer,
        })
    }

    err = fh.TemplatedWriter.Render(httpResponseWriter, "faq_page_template.html", pageObject)
    if err != nil {
        logger.LogError("error rendering faq page" +
                        "|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
    }
}

type PageObject struct {
    FaqItems []Item
}

type Item struct {
    Id string
    Question string
    Answer string
}
