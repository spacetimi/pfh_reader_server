package user_preferences

import (
	"github.com/spacetimi/pfh_reader_server/app_src/app_types"
	"strings"
)

type UserPreferencesData struct {
	CategoryRules []CategoryRule
}

type CategoryRule struct {
	MatchExpression string
	MatchType  CategoryRuleMatchType_t
	IgnoreCase bool

	Category app_types.Category_t
}

type CategoryRuleMatchType_t int
const (
	MATCH_APP_NAME CategoryRuleMatchType_t = iota
	MATCH_APP_TITLE_BAR
)

func (cr *CategoryRule) DoesMatch(appName string, appTitleBar string) bool {
	var s1 string
	s2 := cr.MatchExpression

	switch cr.MatchType {
	case MATCH_APP_NAME: s1 = appName
	case MATCH_APP_TITLE_BAR: s1 = appTitleBar
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