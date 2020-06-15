package app_core

import (
	"os/user"

	"github.com/spacetimi/timi_shared_server/utils/logger"
)

const USER_PREFERENCES_FILE_NAME = "UserPreferencesData.json"

var PFH_DAEMON_DATA_PATH string

/** Package init **/
func init() {
	usr, err := user.Current()
	if err != nil {
		logger.LogFatal("unable to get current user|error=" + err.Error())
	}

	PFH_DAEMON_DATA_PATH = usr.HomeDir + "/Library/Containers/com.spacetimi.pfh-daemon/Data/Documents"
}
