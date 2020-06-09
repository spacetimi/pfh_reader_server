package app_simple_message_page

import (
    "github.com/spacetimi/timi_shared_server/code/config"
    "github.com/spacetimi/timi_shared_server/utils/logger"
    "github.com/spacetimi/timi_shared_server/utils/templated_writer"
    "net/http"
)

func ShowAppSimpleMessagePage(httpResponseWriter http.ResponseWriter,
                             messageHeader string,
                             messageBody string,
                             backlinkHref string,
                             backlinkHrefName string) {

    pageObject := &AppSimpleMessagePageObject{
        MessageHeader:messageHeader,
        MessageBody:messageBody,
        BackLinkHref:backlinkHref,
        BackLinkHrefName:backlinkHrefName,
    }

    page := newAppSimpleMessagePage()
    err := page.Render(httpResponseWriter,
          "app_simple_message_page_template.html",
                       pageObject)
    if err != nil {
        logger.LogError("error showing login message page|error=" + err.Error())
        httpResponseWriter.WriteHeader(http.StatusInternalServerError)
    }
}

////////////////////////////////////////////////////////////////////////////////

type AppSimpleMessagePage struct {
    *templated_writer.TemplatedWriter
}

func newAppSimpleMessagePage() *AppSimpleMessagePage {
    almp := &AppSimpleMessagePage{}
    almp.TemplatedWriter = templated_writer.NewTemplatedWriter(config.GetAppTemplateFilesPath() + "/simple_message_page")

    // Parse templates for every request on LOCAL so that we can iterate over the templates
    // without having to restart the server every time
    almp.TemplatedWriter.ForceReparseTemplates = config.GetEnvironmentConfiguration().AppEnvironment == config.LOCAL

    return almp
}

type AppSimpleMessagePageObject struct {
    MessageHeader string
    MessageBody string
    BackLinkHref string
    BackLinkHrefName string
}
