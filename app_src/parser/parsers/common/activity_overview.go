package common

import (
	"strconv"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/parser_metadata"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"github.com/spacetimi/timi_shared_server/utils/time_utils"
)

type ActivityOverviewData struct {
	ActivityPeriods [parser_metadata.NUM_ACTIVITY_PERIODS_PER_DAY]*ActivityPeriodData
}

type ActivityPeriodData struct {
	PeriodIndex       int
	SecondsInCategory map[app_core.Category_t]int64
}

////////////////////////////////////////////////////////////////////////////////

func NewActivityOverviewData() *ActivityOverviewData {

	aod := &ActivityOverviewData{}

	for i := 0; i < len(aod.ActivityPeriods); i = i + 1 {
		aod.ActivityPeriods[i] = NewActivityPeriodData(i)
	}

	return aod
}

func NewActivityPeriodData(periodIndex int) *ActivityPeriodData {
	apd := &ActivityPeriodData{
		PeriodIndex:       periodIndex,
		SecondsInCategory: make(map[app_core.Category_t]int64),
	}

	for i := app_core.CATEGORY_PRODUCTIVE; i <= app_core.CATEGORY_UNCLASSIFIED; i = i + 1 {
		apd.SecondsInCategory[i] = 0
	}

	return apd
}

func (aod *ActivityOverviewData) AddActivity(category app_core.Category_t, timestampOfActivity int64, seconds int64) {
	timeOfActivity := time_utils.UnixTimeStampToLocalTime(timestampOfActivity)
	beginningOfDay := time_utils.BeginningOfDay(timeOfActivity)

	secondsIntoDay := time_utils.DurationToSeconds(time_utils.GetDurationBetweenTimes(beginningOfDay, timeOfActivity))
	periodIndex := int(secondsIntoDay / parser_metadata.ACTIVITY_PERIOD_LENGTH_SECONDS)

	if periodIndex < 0 || periodIndex >= len(aod.ActivityPeriods) {
		logger.LogError("period index out of range" +
			"|period index=" + strconv.Itoa(periodIndex) +
			"|timestamp of activity=" + strconv.FormatInt(timestampOfActivity, 10))
		return
	}
	activityPeriod := aod.ActivityPeriods[periodIndex]
	activityPeriod.AddActivity(category, seconds)
}

func (apd *ActivityPeriodData) AddActivity(category app_core.Category_t, seconds int64) {
	secondsInCategory, ok := apd.SecondsInCategory[category]
	if !ok {
		secondsInCategory = 0
		apd.SecondsInCategory[category] = secondsInCategory
	}

	apd.SecondsInCategory[category] = secondsInCategory + seconds
}

func (apd *ActivityPeriodData) GetSecondsInCategory(category app_core.Category_t) int64 {
	secondsInCategory, ok := apd.SecondsInCategory[category]
	if !ok {
		return 0
	}
	return secondsInCategory
}

func (aod *ActivityOverviewData) GetActivityInPeriodsForCategory(category app_core.Category_t) []int64 {
	result := make([]int64, len(aod.ActivityPeriods))

	for i, apd := range aod.ActivityPeriods {
		result[i] = apd.GetSecondsInCategory(category)
	}

	return result
}
