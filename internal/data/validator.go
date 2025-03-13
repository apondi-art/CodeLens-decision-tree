package data

import "CodeLens-decision-tree/internal/model"

// ValidateInputFile checks if the input file exists and is a valid CSV
func ValidateInputFile(path string) error

// ValidateColumnExists checks if the specified column exists in the dataset
func ValidateColumnExists(headers []string, columnName string) error

// ValidateConsistentColumns checks if prediction dataset has same columns as training dataset
func ValidateConsistentColumns(trainingHeaders, predictionHeaders []string) error

// ValidateTargetColumn ensures target column data is valid for classification
func ValidateTargetColumn(dataset *model.Dataset, targetColumn string) error

// ValidateDataCompleteness checks for minimum required data to build a model
func ValidateDataCompleteness(dataset *model.Dataset) error

// ValidateDataTypes checks that data types are consistent within columns
func ValidateDataTypes(dataset *model.Dataset) (map[string]string, error)

// ValidateOutputPath ensures the output path is writable
func ValidateOutputPath(path string) error

// ValidateModelFile checks if the model file exists and contains valid JSON
func ValidateModelFile(path string) error
