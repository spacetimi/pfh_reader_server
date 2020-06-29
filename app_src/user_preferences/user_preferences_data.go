package user_preferences

import (
	"strings"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
)

type UserPreferencesData struct {
	CategoryRules []CategoryRule
}

type CategoryRule struct {
	RuleId int

	MatchExpression string
	MatchType       CategoryRuleMatchType_t
	IgnoreCase      bool

	Category app_core.Category_t
}

type CategoryRuleMatchType_t int

const (
	MATCH_APP_NAME CategoryRuleMatchType_t = iota
	MATCH_APP_TITLE_BAR
)

func (usp *UserPreferencesData) GetMatchingCategory(appName string, appTitleBar string) app_core.Category_t {
	for _, cr := range usp.CategoryRules {
		if cr.doesMatch(appName, appTitleBar) {
			return cr.Category
		}
	}

	return app_core.CATEGORY_UNCLASSIFIED
}

func (cr CategoryRule) doesMatch(appName string, appTitleBar string) bool {
	var s1 string
	s2 := cr.MatchExpression

	switch cr.MatchType {
	case MATCH_APP_NAME:
		s1 = appName
	case MATCH_APP_TITLE_BAR:
		s1 = appTitleBar
	default:
		return false
	}

	if cr.IgnoreCase {
		s1 = strings.ToLower(s1)
		s2 = strings.ToLower(s2)
	}

	// TODO: Support match-expression as regex
	return strings.Contains(s1, s2)
}
