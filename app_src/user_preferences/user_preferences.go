package user_preferences

import (
	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/timi_shared_server/utils/file_utils"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

type UserPreferences struct {
	Data *UserPreferencesData
}

var _instance *UserPreferences

func Instance() *UserPreferences {
	if _instance == nil {
		logger.LogFatal("user-preferences instance not set")
	}
	return _instance
}

func CreateInstance() {
	if _instance != nil {
		logger.LogWarning("user-preferences instance already exists. overwriting")
	}

	newInstance := &UserPreferences{}
	userPreferencesDataPath := app_core.PFH_DAEMON_DATA_PATH + "/" + app_core.USER_PREFERENCES_FILE_NAME
	if !file_utils.DoesFileOrDirectoryExist(userPreferencesDataPath) {
		logger.LogFatal("unable to find user-preferences file in path: " + userPreferencesDataPath)
	}

	err := file_utils.ReadJsonFileIntoJsonObject(userPreferencesDataPath, &newInstance.Data)
	if err != nil {
		logger.LogFatal("error reading user-preferences file" +
			"|file path=" + userPreferencesDataPath +
			"|error=" + err.Error())
	}

	_instance = newInstance
}
