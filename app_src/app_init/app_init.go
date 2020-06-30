package app_init

import (
	"errors"

	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/collate"
	"github.com/spacetimi/pfh_reader_server/app_src/user_preferences"
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

	err := CreateBootstrapData()
	if err != nil {
		return errors.New("error creating bootstrap data: " + err.Error())
	}
	user_preferences.CreateInstance()

	collate.CollateDaysToWeeks()

	return nil
}

/********** End IAppInitializer implementation **********/
