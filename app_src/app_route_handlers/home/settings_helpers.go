package home

import (
	"strconv"

	"github.com/spacetimi/pfh_reader_server/app_src/user_preferences"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"github.com/spacetimi/timi_shared_server/utils/slice_utils"
)

func (hh *HomeHandler) getSettingsPageObject() *SettingsData {

	var pageObject *SettingsData
	var appNameMatchRules []SettingsMatchRule
	var appTitleBarMatchRules []SettingsMatchRule

	for _, categoryRule := range user_preferences.Instance().Data.CategoryRules {
		matchRule := SettingsMatchRule{
			RuleId: categoryRule.RuleId,

			MatchExpression: categoryRule.MatchExpression,
			Category:        categoryRule.Category,
			ShouldMatchCase: !categoryRule.IgnoreCase,
		}

		if categoryRule.MatchType == user_preferences.MATCH_APP_NAME {
			appNameMatchRules = append(appNameMatchRules, matchRule)
		} else {
			appTitleBarMatchRules = append(appTitleBarMatchRules, matchRule)
		}
	}

	pageObject = &SettingsData{
		AppNameMatchRules:     appNameMatchRules,
		AppTitleBarMatchRules: appTitleBarMatchRules,
	}
	return pageObject
}

func (hh *HomeHandler) deleteRule(ruleId int) {
	preferencesData := user_preferences.Instance().Data
	if preferencesData == nil {
		logger.LogError("no user preferences data loaded")
		return
	}

	indexInSlice := slice_utils.FindIndexInSlice(len(preferencesData.CategoryRules),
		func(index int) bool {
			return preferencesData.CategoryRules[index].RuleId == ruleId
		})

	if indexInSlice < 0 {
		logger.LogError("unable to find rule to delete" +
			"|rule id to delete=" + strconv.Itoa(ruleId))
		return
	}

	preferencesData.CategoryRules = append(preferencesData.CategoryRules[:indexInSlice],
		preferencesData.CategoryRules[indexInSlice+1:]...)

	err := user_preferences.Instance().SaveChanges()
	if err != nil {
		logger.LogError("error saving user preferences changes while deleting rule" +
			"|rule id to delete=" + strconv.Itoa(ruleId) +
			"|error=" + err.Error())
	}
}
