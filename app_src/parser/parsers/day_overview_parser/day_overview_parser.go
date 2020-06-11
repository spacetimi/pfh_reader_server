package day_overview_parser

import (
	"bufio"
	"errors"
	"github.com/spacetimi/pfh_reader_server/app_src/parser/parsers/parser_metadata"
	"github.com/spacetimi/timi_shared_server/utils/logger"
	"os"
	"strconv"
	"strings"
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

		// TODO: Resolve category from app-name and app-title-bar if override-category is not set
		category := parser_metadata.CATEGORY_UNCLASSIFIED
		if parsedLine.HasOverrideCategory {
			category = parsedLine.OverrideCategory
		}
		dod.AddAppUsageSecondsInCategory(category, parsedLine.AppName, parser_metadata.DAY_LOG_ENTRIES_INTERVAL_SECONDS)
	}

	if err = scanner.Err(); err != nil {
		return nil, errors.New("error parsing file: " + err.Error())
	}

	return dod, nil
}

type parsedLine_t struct {
	Timestamp int64
	AppName string
	AppTitleBar string

	HasOverrideCategory bool
	OverrideCategory parser_metadata.Category_t
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

	overrideCategory := parser_metadata.CATEGORY_UNCLASSIFIED
	hasOverrideCategory := false
	if len(tokens) == 4 {
		hasOverrideCategory = true
		overrideCategory = parser_metadata.CategoryFromString(tokens[3])
	}

	return &parsedLine_t{
		Timestamp: timestamp,
		AppName: appName,
		AppTitleBar: appTitleBar,
		HasOverrideCategory: hasOverrideCategory,
		OverrideCategory: overrideCategory,
	}, nil
}
