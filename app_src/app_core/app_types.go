package app_core

import "strings"

type Category_t int

const (
	CATEGORY_PRODUCTIVE Category_t = iota
	CATEGORY_OPERATIONAL_OVERHEAD
	CATEGORY_UNPRODUCTIVE
	CATEGORY_UNCLASSIFIED
)

func CategoryFromString(s string) (Category_t, bool) {
	switch strings.ToLower(s) {
	case "productive":
		return CATEGORY_PRODUCTIVE, true
	case "operational-overhead":
		return CATEGORY_OPERATIONAL_OVERHEAD, true
	case "unproductive":
		return CATEGORY_UNPRODUCTIVE, true
	case "unclassified":
		return CATEGORY_UNCLASSIFIED, true
	default:
		return -1, false
	}
}
