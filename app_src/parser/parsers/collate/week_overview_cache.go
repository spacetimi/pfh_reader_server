package collate

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"

	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/common"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/week_overview_parser"
	"github.com/spacetimi/timi_shared_server/utils/file_utils"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

var _weekOverviewCacheInstance *WeekOverviewCache

type WeekOverviewCache struct {
	WeekOverviews map[common.WeekIdentifier]*week_overview_parser.WeekOverviewData
}

func WeekOverviewCacheInstance() *WeekOverviewCache {
	if _weekOverviewCacheInstance == nil {
		_weekOverviewCacheInstance = &WeekOverviewCache{
			WeekOverviews: make(map[common.WeekIdentifier]*week_overview_parser.WeekOverviewData),
		}
	}
	return _weekOverviewCacheInstance
}

func (woc *WeekOverviewCache) GetWeekOverview(wid common.WeekIdentifier) (*week_overview_parser.WeekOverviewData, error) {
	wod, ok := woc.WeekOverviews[wid]
	if !ok {
		// week-overview-data hasn't been read or doesn't exist yet

		weekFilePath := common.GetWeekDataFilePath(wid)
		if file_utils.DoesFileOrDirectoryExist(weekFilePath) {
			wod = &week_overview_parser.WeekOverviewData{}
			err := file_utils.ReadJsonFileIntoJsonObject(weekFilePath, wod)
			if err != nil {
				return nil, errors.New("error reading week data file: " + err.Error())
			}
		} else {
			wod = &week_overview_parser.WeekOverviewData{
				WeekIdentifier:        wid,
				WeekdaySummariesByDay: make(map[week_overview_parser.DayOfWeek]*week_overview_parser.WeekdaySummaryData),
			}
		}
		woc.WeekOverviews[wid] = wod
	}

	return wod, nil
}

func (woc *WeekOverviewCache) Apply() {
	for _, wod := range woc.WeekOverviews {
		weekFilePath := common.GetWeekDataFilePath(wod.WeekIdentifier)
		toJson, err := json.MarshalIndent(wod, "", "  ")
		if err != nil {
			logger.LogError("error serializing week data to json" +
				"|week number=" + strconv.Itoa(wod.WeekIdentifier.WeekNumber) +
				"|year number=" + strconv.Itoa(wod.WeekIdentifier.YearNumber) +
				"|error=" + err.Error())
			continue
		}
		err = ioutil.WriteFile(weekFilePath, toJson, 0644)
		if err != nil {
			logger.LogError("error writing week data file" +
				"|file path=" + weekFilePath +
				"|error=" + err.Error())
			continue
		}
	}
	_weekOverviewCacheInstance = nil
}
