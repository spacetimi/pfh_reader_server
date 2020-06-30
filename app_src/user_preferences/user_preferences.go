package user_preferences

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sort"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/timi_shared_server/utils/file_utils"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

type UserPreferences struct {
	DataModTime int64
	Data        *UserPreferencesData
}

var _instance *UserPreferences

func Instance() *UserPreferences {
	if _instance == nil ||
		MustCreateBootstrapData() {
		CreateInstance()
	}
	return _instance
}

func CreateInstance() {

	newInstance := &UserPreferences{}
	userPreferencesDataPath := getUserPreferencesDataFilePath()

	err := CheckAndCreateBootstrapData()
	if err != nil {
		logger.LogError("error creating bootstrap data" +
			"|file path=" + userPreferencesDataPath +
			"|error=" + err.Error())
		return
	}

	err = file_utils.ReadJsonFileIntoJsonObject(userPreferencesDataPath, &newInstance.Data)
	if err != nil {
		logger.LogError("error reading user-preferences file" +
			"|file path=" + userPreferencesDataPath +
			"|error=" + err.Error())
		return
	}

	newInstance.DataModTime, err = file_utils.GetFileModTimeUnix(userPreferencesDataPath)
	if err != nil {
		logger.LogError("error getting user-preferences file mod time" +
			"|file path=" + userPreferencesDataPath +
			"|error=" + err.Error())
		return
	}

	// Sort the category rules so that app-title-bar-rules are checked before app-name-rules
	sort.Slice(newInstance.Data.CategoryRules, func(i, j int) bool {
		return newInstance.Data.CategoryRules[i].MatchType > newInstance.Data.CategoryRules[j].MatchType
	})

	_instance = newInstance
}

func (up *UserPreferences) SaveChanges() error {
	userPreferencesDataPath := app_core.PFH_DAEMON_DATA_PATH + "/" + app_core.USER_PREFERENCES_FILE_NAME
	toJson, err := json.MarshalIndent(up.Data, "", "  ")
	if err != nil {
		return errors.New("error serializing user preferences data to json: " + err.Error())
	}
	err = ioutil.WriteFile(userPreferencesDataPath, toJson, 0644)
	if err != nil {
		return errors.New("error saving user preferences data file: " + err.Error())
	}

	return nil
}

func Reload() {
	if _instance != nil {
		_instance = nil
	}
	CreateInstance()
}

func getUserPreferencesDataFilePath() string {
	return app_core.PFH_DAEMON_DATA_PATH + "/" + app_core.USER_PREFERENCES_FILE_NAME
}
