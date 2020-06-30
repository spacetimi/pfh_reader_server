package user_preferences

import (
	"errors"

	"github.com/spacetimi/timi_shared_server/code/config"
	"github.com/spacetimi/timi_shared_server/utils/file_utils"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

const kBootstrapFileName = "UserPreferencesBootstrapData.json"

func CheckAndCreateBootstrapData() error {

	bootstrapFilePath := getBootstrapFilepath()
	targetFilePath := getUserPreferencesDataFilePath()

	if !file_utils.DoesFileOrDirectoryExist(bootstrapFilePath) {
		return errors.New("no bootstrap file present in path: " + bootstrapFilePath)
	}

	if !MustCreateBootstrapData() {
		return nil
	}

	err := file_utils.CopyFile(bootstrapFilePath, targetFilePath)
	if err != nil {
		logger.LogError("error copying bootstrap file" +
			"|source file=" + bootstrapFilePath +
			"|target file=" + targetFilePath +
			"|error=" + err.Error())
		return errors.New("error copying bootstrap file: " + err.Error())
	}

	logger.LogInfo("successfully copied bootstrap data")

	return nil
}

/*
 Do we need to copy over bootstrap data?
 Check if:
     1. the target file is missing
     2. or it is out of data (which could happen if there is a newer release of pfh-reader which has an updated bootstrap file)
*/
func MustCreateBootstrapData() bool {
	bootstrapFilepath := getBootstrapFilepath()
	targetFilepath := getUserPreferencesDataFilePath()

	if !file_utils.DoesFileOrDirectoryExist(targetFilepath) {
		return true
	}

	targetFileModTime, err := file_utils.GetFileModTimeUnix(targetFilepath)
	if err != nil {
		return true
	}

	bootstrapFileModTime, err := file_utils.GetFileModTimeUnix(bootstrapFilepath)
	if err != nil {
		return true
	}

	return bootstrapFileModTime > targetFileModTime
}

func getBootstrapFilepath() string {
	return config.GetAppResourcesPath() + "/bootstrap_data/" + kBootstrapFileName
}
