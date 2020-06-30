package user_preferences

import (
	"errors"

	"github.com/spacetimi/timi_shared_server/code/config"
	"github.com/spacetimi/timi_shared_server/utils/file_utils"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

const kBootstrapFileName = "UserPreferencesBootstrapData.json"

func CreateBootstrapData(targetFilePath string) error {

	if file_utils.DoesFileOrDirectoryExist(targetFilePath) {
		return nil
	}

	bootstrapFilePath := config.GetAppResourcesPath() + "/bootstrap_data/" + kBootstrapFileName
	if !file_utils.DoesFileOrDirectoryExist(bootstrapFilePath) {
		return errors.New("no bootstrap file present in path: " + bootstrapFilePath)
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
