package home

import (
	"strconv"

	"github.com/spacetimi/timi_shared_server/utils/logger"
)

type HomePageTab_t int

const (
	HOMEPAGE_TAB_DASHBOARD HomePageTab_t = iota
	HOMEPAGE_TAB_WEEK
	HOMEPAGE_TAB_INSIGHTS
	HOMEPAGE_TAB_SETTINGS
)

const HOMEPAGE_TABNAME_DASHBOARD = "dashboard"
const HOMEPAGE_TABNAME_WEEK = "week"
const HOMEPAGE_TABNAME_INSIGHTS = "insights"
const HOMEPAGE_TABNAME_SETTINGS = "settings"

func parseHomePageTabString(s string) HomePageTab_t {
	switch s {
	case HOMEPAGE_TABNAME_DASHBOARD:
		return HOMEPAGE_TAB_DASHBOARD
	case HOMEPAGE_TABNAME_WEEK:
		return HOMEPAGE_TAB_WEEK
	case HOMEPAGE_TABNAME_INSIGHTS:
		return HOMEPAGE_TAB_INSIGHTS
	case HOMEPAGE_TABNAME_SETTINGS:
		return HOMEPAGE_TAB_SETTINGS
	}
	return HOMEPAGE_TAB_DASHBOARD
}

func (tab HomePageTab_t) String() string {
	switch tab {
	case HOMEPAGE_TAB_DASHBOARD:
		return HOMEPAGE_TABNAME_DASHBOARD
	case HOMEPAGE_TAB_WEEK:
		return HOMEPAGE_TABNAME_WEEK
	case HOMEPAGE_TAB_INSIGHTS:
		return HOMEPAGE_TABNAME_INSIGHTS
	case HOMEPAGE_TAB_SETTINGS:
		return HOMEPAGE_TABNAME_SETTINGS
	}
	return HOMEPAGE_TABNAME_DASHBOARD
}

////////////////////////////////////////////////////////////////////////////////

const kPostArgNameCurrentTab = "tab"
const kPostArgNameCurrentDayIndex = "day-index"

type parsedPostArgs struct {
	Tab             HomePageTab_t
	CurrentDayIndex int // 0 is today, -1 is yesterday, and so on
}

func parsePostArgs(postArgs map[string]string) *parsedPostArgs {
	parsed := &parsedPostArgs{
		Tab:             HOMEPAGE_TAB_DASHBOARD,
		CurrentDayIndex: 0,
	}

	if postArgs == nil || len(postArgs) == 0 {
		return parsed
	}

	currentTabName, ok := postArgs[kPostArgNameCurrentTab]
	if ok {
		parsed.Tab = parseHomePageTabString(currentTabName)
	}

	dayIndexString, ok := postArgs[kPostArgNameCurrentDayIndex]
	if ok {
		dayIndex, err := strconv.ParseInt(dayIndexString, 10, 32)
		if err != nil {
			logger.LogError("error parsing day index from post args" +
				"|day index string=" + dayIndexString +
				"|error=" + err.Error())
		} else {
			parsed.CurrentDayIndex = int(dayIndex)
		}
	}

	return parsed
}