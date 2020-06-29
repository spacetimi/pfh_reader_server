package home

import "github.com/spacetimi/pfh_reader_server/app_src/user_preferences"

func (hh *HomeHandler) getSettingsPageObject() *SettingsData {

	var pageObject *SettingsData
	var matchRules []SettingsMatchRule

	for _, categoryRule := range user_preferences.Instance().Data.CategoryRules {
		matchRule := SettingsMatchRule{
			MatchType:       categoryRule.MatchType,
			MatchExpression: categoryRule.MatchExpression,
			Category:        categoryRule.Category,
			ShouldMatchCase: !categoryRule.IgnoreCase,
		}
		matchRules = append(matchRules, matchRule)
	}

	pageObject = &SettingsData{
		MatchRules: matchRules,
	}
	return pageObject
}
