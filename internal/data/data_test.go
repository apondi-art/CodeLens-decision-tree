package data

import (
	"os"
	"testing"

	"CodeLens-decision-tree/internal/model"
)

func TestParseNumerical(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		wantErr  bool
	}{
		{"42.5", 42.5, false},
		{"-10.3", -10.3, false},
		{"0", 0.0, false},
		{"", 0.0, true},       // Empty value should return an error
		{"abc", 0.0, true},    // Invalid numerical input
		{"12.3.4", 0.0, true}, // Malformed number
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseNumerical(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if result != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestParseCategorical(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"Apple", "Apple", false},
		{"  Banana  ", "Banana", false}, // Trimming whitespace
		{"", "", true},                  // Empty value should return an error
		{" ", "", true},                 // Whitespace-only should return an error
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseCategorical(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			if result != tt.expected {
				t.Errorf("expected: %q, got: %q", tt.expected, result)
			}
		})
	}
}

func TestParseTimestamp(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{"Valid YYYY-MM-DD", "2025-03-14", "2025-03-14 00:00:00", false},
		{"Valid YYYY-MM-DD HH:MM:SS", "2025-03-14 12:30:45", "2025-03-14 12:30:45", false},
		{"Valid ISO 8601", "2025-03-14T12:30:45Z", "2025-03-14 12:30:45", false},
		{"Valid DD-MM-YYYY", "14-03-2025", "2025-03-14 00:00:00", false},
		{"Invalid format", "invalid-date", "", true},
		{"Non-existent date", "2025-99-99", "", true},
		{"Empty input", "", "", true},
		{"Only time", "12:30:45", "", true},
		{"Future date", "3025-12-25", "3025-12-25 00:00:00", false},
		{"Oldest valid date", "0001-01-01", "0001-01-01 00:00:00", false},
		{"Leap year date", "2024-02-29", "2024-02-29 00:00:00", false},
		{"Non-leap year invalid date", "2023-02-29", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTimestamp(tt.input)

			// Error check
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}

			// Compare valid results
			if !tt.wantErr {
				resultStr := result.Format("2006-01-02 15:04:05")
				if resultStr != tt.expected {
					t.Errorf("expected: %v, got: %v", tt.expected, resultStr)
				}
			}
		})
	}
}

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
