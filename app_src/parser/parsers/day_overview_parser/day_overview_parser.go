package day_overview_parser

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/spacetimi/pfh_reader_server/app_src/app_core"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/parser_metadata"
	"github.com/spacetimi/pfh_reader_server/app_src/user_preferences"
	"github.com/spacetimi/timi_shared_server/utils/logger"
)

type DayOverviewParser struct {
}

func (dop *DayOverviewParser) ParseFile(filePath string) (*DayOverviewData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("error opening file: " + err.Error())
	}

	defer func() {
		err := file.Close()
		if err != nil {
			logger.LogError("error closing day-file after parsing" +
				"|file path=" + filePath +
				"|error=" + err.Error())
		}
	}()

	dod := NewDayOverviewData()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parsedLine, parse_err := parseLine(line)
		if parse_err != nil {
			continue
		}
		if parsedLine.AppName == "locked" {
			continue
		}

		category := app_core.CATEGORY_UNCLASSIFIED
		if parsedLine.HasOverrideCategory {
			category = parsedLine.OverrideCategory
		} else {
			category = user_preferences.Instance().Data.GetMatchingCategory(parsedLine.AppName, parsedLine.AppTitleBar)
		}
		dod.AddAppUsageSecondsInCategory(category, parsedLine.AppName, parsedLine.Timestamp, parser_metadata.DAY_LOG_ENTRIES_INTERVAL_SECONDS)
	}

	if err = scanner.Err(); err != nil {
		return nil, errors.New("error parsing file: " + err.Error())
	}

	return dod, nil
}

type parsedLine_t struct {
	Timestamp   int64
	AppName     string
	AppTitleBar string

	HasOverrideCategory bool
	OverrideCategory    app_core.Category_t
}

func parseLine(line string) (*parsedLine_t, error) {
	tokens := strings.Split(line, parser_metadata.DAY_LOG_ENTRY_SEPARATOR)
	if len(tokens) < 3 || len(tokens) > 4 {
		return nil, errors.New("unexpected number of tokens found while parsing line: " + line)
	}

	timestamp, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return nil, errors.New("error parsing timestamp in line: " + line)
	}
	appName := tokens[1]
	appTitleBar := tokens[2]

	overrideCategory := app_core.CATEGORY_UNCLASSIFIED
	hasOverrideCategory := false
	if len(tokens) == 4 {
		overrideCategory, hasOverrideCategory = app_core.CategoryFromString(tokens[3])
	}

	return &parsedLine_t{
		Timestamp:           timestamp,
		AppName:             appName,
		AppTitleBar:         appTitleBar,
		HasOverrideCategory: hasOverrideCategory,
		OverrideCategory:    overrideCategory,
	}, nil
}
