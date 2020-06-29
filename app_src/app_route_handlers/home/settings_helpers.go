package home

import "github.com/spacetimi/pfh_reader_server/app_src/user_preferences"

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
