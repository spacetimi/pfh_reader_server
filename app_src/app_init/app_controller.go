package app_init

import (
    "github.com/spacetimi/pfh_reader_server/app_src/app_route_handlers/about"
    "github.com/spacetimi/pfh_reader_server/app_src/app_route_handlers/faq"
    "github.com/spacetimi/pfh_reader_server/app_src/app_route_handlers/home"
    "github.com/spacetimi/timi_shared_server/code/core/controller"
)

type AppController struct { // Implements IAppController
}

func (ac *AppController) RouteHandlers() []controller.IRouteHandler {
    return []controller.IRouteHandler {
        home.NewHomeHandler(),
        faq.NewFaqHandler(),
        about.NewAboutHandler(),
    }
}

