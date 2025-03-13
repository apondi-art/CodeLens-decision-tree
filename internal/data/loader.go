package data

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"CodeLens-decision-tree/internal/model"
)



func GenerateCSVData(path, targetColumn string) (*model.Dataset, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dataset := &model.Dataset{
		RowInstances:    []map[string]interface{}{},
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
			Max:           -1e9,
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
	if targetIndex == -1 {
		return nil, fmt.Errorf("Target column '%s' not found in CSV", targetColumn)
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

			if numericVal, err := strconv.ParseFloat(value, 64); err == nil {
				rowInstance[colName] = numericVal
				attr.Type = model.Numeric

				if numericVal < attr.Min {
					attr.Min = numericVal
				}
				if numericVal > attr.Max {
					attr.Max = numericVal
				}
			} else {
				rowInstance[colName] = value
				attr.Type = model.Categorical

				// Track unique values using a map for efficiency
				uniqueValues := make(map[string]struct{}, len(attr.PossibleValues))
				for _, v := range attr.PossibleValues {
					uniqueValues[v.(string)] = struct{}{}
				}

				if _, exists := uniqueValues[value]; !exists {
					attr.PossibleValues = append(attr.PossibleValues, value)
				}
			}
		}

		dataset.RowInstances = append(dataset.RowInstances, rowInstance)

		// Update TargetOccurrence
		if classValue, ok := rowInstance[targetColumn].(string); ok {
			dataset.TargetOccurrence[classValue]++
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
// func InferDataTypes(data [][]string) (map[string]string, error)
// func HandleMissingValues(dataset *model.Dataset) error
//
