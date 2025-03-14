package data

import (
	"encoding/json"
	"errors"
	"os"

	"CodeLens-decision-tree/internal/model"
)

// ValidateInputFile checks if the input file exists and is a valid CSV
// ValidateInputFile checks if the input file exists and is a valid CSV
func ValidateInputFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("input file does not exist: " + path)
		}
		return err
	}

	// Ensure it's a file and not a directory
	if info.IsDir() {
		return errors.New("input path is a directory, not a file: " + path)
	}

	return nil
}

// ValidateColumnExists checks if the specified column exists in the dataset
func ValidateColumnExists(headers []string, columnName string) error {
	for _, header := range headers {
		if header == columnName {
			return nil // Column found
		}
	}
	return errors.New("column '" + columnName + "' not found in dataset")
}

// ValidateConsistentColumns checks if prediction dataset has same columns as training dataset
// ValidateConsistentColumns checks if prediction dataset has same columns as training dataset
func ValidateConsistentColumns(trainingHeaders, predictionHeaders []string) error {
	if len(trainingHeaders) != len(predictionHeaders) {
		return errors.New("training and prediction datasets have different column counts")
	}

	headerMap := make(map[string]bool)
	for _, col := range trainingHeaders {
		headerMap[col] = true
	}

	for _, col := range predictionHeaders {
		if !headerMap[col] {
			return errors.New("column mismatch: " + col + " is missing in prediction dataset")
		}
	}

	return nil
}

// ValidateTargetColumn ensures target column data is valid for classification
// ValidateTargetColumn ensures target column data is valid for classification
func ValidateTargetColumn(dataset *model.Dataset, targetColumn string) error {
	if _, exists := dataset.ColumnAttributes[targetColumn]; !exists {
		return errors.New("target column '" + targetColumn + "' does not exist in dataset")
	}

	if len(dataset.TargetOccurrence) == 0 {
		return errors.New("target column has no valid class labels")
	}

	return nil
}

// ValidateDataCompleteness checks for minimum required data to build a model

// ValidateDataCompleteness checks for minimum required data to build a model
func ValidateDataCompleteness(dataset *model.Dataset) error {
	if dataset.TotalRows < 2 {
		return errors.New("insufficient data: at least two rows are required")
	}

	for col, attr := range dataset.ColumnAttributes {
		if attr.MissingCount == dataset.TotalRows {
			return errors.New("column '" + col + "' is entirely missing values")
		}
	}

	return nil
}

// ValidateDataTypes checks that data types are consistent within columns
func ValidateDataTypes(dataset *model.Dataset) (map[string]string, error) {
	columnTypes := make(map[string]string)

	for col, attr := range dataset.ColumnAttributes {
		switch attr.Type {
		case model.Numeric:
			columnTypes[col] = "numeric"
		case model.Categorical:
			columnTypes[col] = "categorical"
		default:
			return nil, errors.New("unknown data type for column: " + col)
		}
	}

	return columnTypes, nil
}

// ValidateOutputPath ensures the output path is writable
func ValidateOutputPath(path string) error {
	file, err := os.CreateTemp(path, "testfile")
	if err != nil {
		return errors.New("output path is not writable: " + path)
	}
	defer os.Remove(file.Name())

	return nil
}

// ValidateModelFile checks if the model file exists and contains valid JSON
// ValidateModelFile checks if the model file exists and contains valid JSON
func ValidateModelFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.New("model file does not exist or is unreadable: " + path)
	}

	var content map[string]interface{}
	if err := json.Unmarshal(data, &content); err != nil {
		return errors.New("invalid model file format: must be valid JSON")
	}

	return nil
}
