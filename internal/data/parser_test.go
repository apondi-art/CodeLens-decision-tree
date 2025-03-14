package data

import (
	"testing"
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
