package data

import "time"

func ParseNumerical(value string) (float64, error) {

	return 0.0, nil
}

func ParseCategorical(value string) (string, error) {
	return "", nil
}

func ParseTimestamp(value string) (time.Time, error) {
	return time.Now(), nil
}
