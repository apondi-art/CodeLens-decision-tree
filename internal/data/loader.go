package data

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"CodeLens-decision-tree/internal/model"
)

// GenerateCSVData reads a CSV file, processes data using parsers, and returns a structured dataset.
func GenerateCSVData(path, targetColumn string) (*model.Dataset, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dataset := &model.Dataset{
		RowInstances:     []map[string]interface{}{},
		ColumnAttributes: make(map[string]*model.Attribute),
		TargetOccurrence: make(map[string]int),
	}

	scanner := bufio.NewScanner(file)

	// Read header row
	if !scanner.Scan() {
		return nil, fmt.Errorf("CSV file is empty")
	}

	dataset.ColumnNames = strings.Split(strings.TrimSpace(scanner.Text()), ",")

	// Initialize attributes
	for _, col := range dataset.ColumnNames {
		dataset.ColumnAttributes[col] = &model.Attribute{
			Name:           col,
			Type:           model.Unknown,
			PossibleValues: []interface{}{},
			Min:            1e9,
			Max:            -1e9,
			MissingCount:   0,
		}
	}

	// Find target column index
	targetIndex := -1
	for i, col := range dataset.ColumnNames {
		if col == targetColumn {
			targetIndex = i
			break
		}
	}
	if targetIndex == -1 && targetColumn != "" {
		return nil, fmt.Errorf("target column '%s' not found in CSV", targetColumn)
	}
	dataset.TargetColumn = targetColumn

	// Read data rows
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		rowValues := strings.Split(line, ",")

		if len(rowValues) < len(dataset.ColumnNames) {
			continue // Skip incomplete rows
		}

		rowInstance := make(map[string]interface{})

		for i, colName := range dataset.ColumnNames {
			value := strings.TrimSpace(rowValues[i])
			attr := dataset.ColumnAttributes[colName]

			if value == "" {
				rowInstance[colName] = nil
				attr.MissingCount++
				continue
			}

			// Try parsing as a numerical value
			if numericVal, err := ParseNumerical(value); err == nil {
				rowInstance[colName] = numericVal
				attr.Type = model.Numeric

				if numericVal < attr.Min {
					attr.Min = numericVal
				}
				if numericVal > attr.Max {
					attr.Max = numericVal
				}
				continue
			}

			// Try parsing as a timestamp
			if timestampVal, err := ParseTimestamp(value); err == nil {
				rowInstance[colName] = timestampVal
				attr.Type = model.Timestamp
				continue
			}

			// Treat as categorical data
			if categoricalVal, err := ParseCategorical(value); err == nil {
				rowInstance[colName] = categoricalVal
				attr.Type = model.Categorical

				// Track unique values using a map for efficiency
				uniqueValues := make(map[string]struct{}, len(attr.PossibleValues))
				for _, v := range attr.PossibleValues {
					uniqueValues[v.(string)] = struct{}{}
				}

				if _, exists := uniqueValues[categoricalVal]; !exists {
					attr.PossibleValues = append(attr.PossibleValues, categoricalVal)
				}
				continue
			}

			// If all parsing attempts fail, mark as missing
			rowInstance[colName] = nil
			attr.MissingCount++
		}

		dataset.RowInstances = append(dataset.RowInstances, rowInstance)

		// Update TargetOccurrence
		if targetColumn != "" {
			if classValue, ok := rowInstance[targetColumn].(string); ok {
				dataset.TargetOccurrence[classValue]++
			}
		}
	}

	// Set metadata values
	dataset.TotalRows = len(dataset.RowInstances)
	dataset.NonTargetColumns = len(dataset.ColumnNames) - 1

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dataset, nil
}
