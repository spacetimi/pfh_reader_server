package app_init

import (
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/collate"
	"github.com/spacetimi/timi_shared_server/code/core/shared_init"
)

func GetAppInitializer() shared_init.IAppInitializer {
	return &appInitializer
}

type AppInitializer struct { // Implements IAppInit
}

var appInitializer AppInitializer

/********** Begin IAppInitializer implementation **********/

func (appInitializer *AppInitializer) AppName() string {
	return "pfh_reader_server"
}

func (appInitializer *AppInitializer) AppInit() error {

	collate.CollateDaysToWeeks()

	return nil
}

/********** End IAppInitializer implementation **********/
