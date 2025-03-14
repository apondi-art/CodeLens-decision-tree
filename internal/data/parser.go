package data

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// ParseNumerical converts a string to a float64.
// Returns an error if the input is empty or not a valid number.
func ParseNumerical(value string) (float64, error) {
	if value == "" {
		return 0.0, errors.New("empty value cannot be converted to number")
	}
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0, errors.New("invalid numerical value: " + value)
	}
	return num, nil
}

// ParseCategorical trims whitespace and validates a categorical string.
// Returns an error if the input is empty.
func ParseCategorical(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.New("empty categorical value")
	}
	return value, nil
}

// ParseTimestamp parses a string into a time.Time object using common date-time formats.
// Returns an error if the input format is invalid.

func ParseTimestamp(value string) (time.Time, error) {
	formats := []string{
		"2006-01-02",                // YYYY-MM-DD
		"2006-01-02 15:04:05",       // YYYY-MM-DD HH:MM:SS
		"2006-01-02T15:04:05Z07:00", // ISO 8601
		"02-01-2006",                // DD-MM-YYYY
	}

	for _, format := range formats {
		parsedTime, err := time.Parse(format, value)
		if err == nil {
			return parsedTime, nil
		}
	}

	// Return zero time explicitly when parsing fails
	return time.Time{}, errors.New("invalid timestamp format: " + value)
}
