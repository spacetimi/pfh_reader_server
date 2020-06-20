package parser_metadata

const DAY_LOG_ENTRIES_INTERVAL_SECONDS = 15
const DAY_LOG_ENTRY_SEPARATOR = "####"

const SECONDS_PER_DAY = 3600 * 24
const ACTIVITY_PERIOD_LENGTH_SECONDS = 60 * 60 // 60 minutes
const NUM_ACTIVITY_PERIODS_PER_DAY = SECONDS_PER_DAY / ACTIVITY_PERIOD_LENGTH_SECONDS

func ParseActivityPeriodIndex(index int) (hours int, minutes int) {
	seconds := index * ACTIVITY_PERIOD_LENGTH_SECONDS
	hours = seconds / 3600
	minutes = (seconds % 3600) / 60

	return hours, minutes
}
