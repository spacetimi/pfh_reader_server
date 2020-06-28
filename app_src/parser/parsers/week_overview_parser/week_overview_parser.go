package week_overview_parser

import (
	"errors"

	"github.com/spacetimi/timi_shared_server/utils/file_utils"
)

type WeekOverviewParser struct {
}

func (wop *WeekOverviewParser) ParseFile(filePath string) (*WeekOverviewData, error) {
	if !file_utils.DoesFileOrDirectoryExist(filePath) {
		return nil, errors.New("no such file")
	}

	wod := &WeekOverviewData{}
	err := file_utils.ReadJsonFileIntoJsonObject(filePath, wod)
	if err != nil {
		return nil, errors.New("error reading file: " + err.Error())
	}

	return wod, nil
}
