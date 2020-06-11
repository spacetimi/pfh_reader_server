package parser_metadata

import "strings"

type Category_t int
const (
	CATEGORY_PRODUCTIVE Category_t = iota
	CATEGORY_OPERATIONAL_OVERHEAD
	CATEGORY_UNPRODUCTIVE
	CATEGORY_UNCLASSIFIED
)

func CategoryFromString(s string) Category_t {
	switch strings.ToLower(s) {
	case "productive": return CATEGORY_PRODUCTIVE
	case "operational_overhead": return CATEGORY_OPERATIONAL_OVERHEAD
	case "unproductive": return CATEGORY_UNPRODUCTIVE
	default: return CATEGORY_UNCLASSIFIED
	}
}

const DAY_LOG_ENTRIES_INTERVAL_SECONDS = 15
const DAY_LOG_ENTRY_SEPARATOR = "####"
