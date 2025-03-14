package data

import (
	"os"
	"testing"

	"CodeLens-decision-tree/internal/model"
)

// TestValidateInputFile - Covers cases where the file exists, does not exist, or is a directory
func TestValidateInputFile(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "testfile.csv")
	if err != nil {
		t.Fatal("Failed to create temp file")
	}
	defer os.Remove(tempFile.Name())

	tests := []struct {
		name      string
		filePath  string
		expectErr bool
	}{
		{"Valid File", tempFile.Name(), false},
		{"Non-existent File", "nonexistent.csv", true},
		{"Directory Path", ".", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInputFile(tt.filePath)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

// TestValidateColumnExists - Covers cases where the column exists and does not exist
func TestValidateColumnExists(t *testing.T) {
	headers := []string{"age", "salary", "status"}

	tests := []struct {
		name       string
		columnName string
		expectErr  bool
	}{
		{"Column Exists", "age", false},
		{"Column Does Not Exist", "gender", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateColumnExists(headers, tt.columnName)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

// TestValidateConsistentColumns - Checks if columns match between datasets
func TestValidateConsistentColumns(t *testing.T) {
	tests := []struct {
		name              string
		trainingHeaders   []string
		predictionHeaders []string
		expectErr         bool
	}{
		{"Matching Columns", []string{"age", "salary"}, []string{"age", "salary"}, false},
		{"Different Column Count", []string{"age", "salary"}, []string{"age"}, true},
		{"Missing Column", []string{"age", "salary"}, []string{"age", "gender"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConsistentColumns(tt.trainingHeaders, tt.predictionHeaders)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

// TestValidateTargetColumn - Tests for existence and validity of target column
func TestValidateTargetColumn(t *testing.T) {
	tests := []struct {
		name      string
		dataset   *model.Dataset
		targetCol string
		expectErr bool
	}{
		{"Valid Target Column", &model.Dataset{
			ColumnAttributes: map[string]*model.Attribute{"status": {}},
			TargetOccurrence: map[string]int{"active": 1}, // Add at least one class label
		}, "status", false},
		{"Missing Target Column", &model.Dataset{ColumnAttributes: map[string]*model.Attribute{}}, "status", true},
		{"Empty Target Occurrence", &model.Dataset{
			ColumnAttributes: map[string]*model.Attribute{"status": {}},
			TargetOccurrence: map[string]int{}, // Empty class labels
		}, "status", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTargetColumn(tt.dataset, tt.targetCol)
			if (err != nil) != tt.expectErr {
				t.Errorf("Test '%s' failed: Expected error: %v, got: %v", tt.name, tt.expectErr, err)
			}
		})
	}
}

// TestValidateDataCompleteness - Tests if dataset has enough valid data
func TestValidateDataCompleteness(t *testing.T) {
	tests := []struct {
		name      string
		dataset   *model.Dataset
		expectErr bool
	}{
		{"Sufficient Data", &model.Dataset{TotalRows: 5, ColumnAttributes: map[string]*model.Attribute{"age": {MissingCount: 0}}}, false},
		{"Insufficient Rows", &model.Dataset{TotalRows: 1}, true},
		{"All Missing Values", &model.Dataset{TotalRows: 5, ColumnAttributes: map[string]*model.Attribute{"age": {MissingCount: 5}}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDataCompleteness(tt.dataset)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

// TestValidateDataTypes - Ensures data types are correct
func TestValidateDataTypes(t *testing.T) {
	tests := []struct {
		name      string
		dataset   *model.Dataset
		expectErr bool
	}{
		{"Valid Data Types", &model.Dataset{ColumnAttributes: map[string]*model.Attribute{"age": {Type: model.Numeric}, "gender": {Type: model.Categorical}}}, false},
		{"Unknown Data Type", &model.Dataset{ColumnAttributes: map[string]*model.Attribute{"age": {Type: -1}}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateDataTypes(tt.dataset)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

// TestValidateOutputPath - Checks if output path is writable
func TestValidateOutputPath(t *testing.T) {
	tempDir := os.TempDir()

	tests := []struct {
		name      string
		path      string
		expectErr bool
	}{
		{"Writable Path", tempDir, false},
		{"Non-existent Path", "/invalid/path", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOutputPath(tt.path)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

// TestValidateModelFile - Ensures model file exists and is valid JSON
func TestValidateModelFile(t *testing.T) {
	// Create a valid JSON model file
	validFile, err := os.CreateTemp("", "valid_model.json")
	if err != nil {
		t.Fatal("Failed to create temp file")
	}
	defer os.Remove(validFile.Name())
	validContent := `{"tree": {"decision": "yes"}}`
	validFile.Write([]byte(validContent))
	validFile.Close()

	// Create an invalid JSON model file
	invalidFile, err := os.CreateTemp("", "invalid_model.json")
	if err != nil {
		t.Fatal("Failed to create temp file")
	}
	defer os.Remove(invalidFile.Name())
	invalidFile.Write([]byte("{invalid json"))
	invalidFile.Close()

	tests := []struct {
		name      string
		filePath  string
		expectErr bool
	}{
		{"Valid JSON Model File", validFile.Name(), false},
		{"Invalid JSON Model File", invalidFile.Name(), true},
		{"Non-existent Model File", "missing_model.json", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateModelFile(tt.filePath)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}
