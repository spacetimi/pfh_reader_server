package home

import (
	"errors"
	"strconv"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/user_preferences"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"github.com/spacetimi/timi_shared_server/utils/slice_utils"
)

func (hh *HomeHandler) getSettingsPageObject() *SettingsData {

	var pageObject *SettingsData
	var appNameMatchRules []SettingsMatchRule
	var appTitleBarMatchRules []SettingsMatchRule

	if user_preferences.Instance() == nil {
		return &SettingsData{
			ErrorablePage: ErrorablePage{
				HasError:    true,
				ErrorString: "Error loading user preferences",
			},
		}
	}

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
		ErrorablePage: ErrorablePage{
			HasError:    false,
			ErrorString: "",
		},
		AppNameMatchRules:     appNameMatchRules,
		AppTitleBarMatchRules: appTitleBarMatchRules,
	}
	return pageObject
}

func (hh *HomeHandler) deleteRule(ruleId int) error {
	if user_preferences.Instance() == nil {
		return errors.New("error getting user preferences")
	}

	preferencesData := user_preferences.Instance().Data
	if preferencesData == nil {
		return errors.New("no user preferences data loaded")
	}

	indexInSlice := slice_utils.FindIndexInSlice(len(preferencesData.CategoryRules),
		func(index int) bool {
			return preferencesData.CategoryRules[index].RuleId == ruleId
		})

	if indexInSlice < 0 {
		logger.LogError("unable to find rule to delete" +
			"|rule id to delete=" + strconv.Itoa(ruleId))
		return errors.New("no such rule")
	}

	preferencesData.CategoryRules = append(preferencesData.CategoryRules[:indexInSlice],
		preferencesData.CategoryRules[indexInSlice+1:]...)

	err := user_preferences.Instance().SaveChanges()
	if err != nil {
		logger.LogError("error saving user preferences changes while deleting rule" +
			"|rule id to delete=" + strconv.Itoa(ruleId) +
			"|error=" + err.Error())
		return errors.New("unable to save user preferences")
	}

	return nil
}

func (hh *HomeHandler) addRule(matchType user_preferences.CategoryRuleMatchType_t,
	matchExpression string, matchCase bool, category app_core.Category_t) error {

	if user_preferences.Instance() == nil {
		return errors.New("error getting user preferences")
	}

	preferencesData := user_preferences.Instance().Data
	if preferencesData == nil {
		return errors.New("no user preferences data loaded")
	}

	highestRuleId := 0
	for _, rule := range preferencesData.CategoryRules {
		if rule.RuleId > highestRuleId {
			highestRuleId = rule.RuleId
		}
	}

	newRule := user_preferences.CategoryRule{
		RuleId:          highestRuleId + 1,
		MatchType:       matchType,
		MatchExpression: matchExpression,
		IgnoreCase:      !matchCase,
		Category:        category,
	}

	preferencesData.CategoryRules = append(preferencesData.CategoryRules, newRule)

	err := user_preferences.Instance().SaveChanges()
	if err != nil {
		logger.LogError("error saving user preferences changes while adding new rule" +
			"|error=" + err.Error())
		return errors.New("error saving user preferences")
	}

	return nil
}
