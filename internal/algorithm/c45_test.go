package algorithm

import (
	"testing"

	"CodeLens-decision-tree/internal/model"
)

func TestMajorityClass(t *testing.T) {
	tests := []struct {
		name     string
		dataset  *model.Dataset
		target   string
		expected string
	}{
		{
			name: "Single class",
			dataset: &model.Dataset{
				RowInstances: []map[string]interface{}{
					{"PlayTennis": "Yes"},
					{"PlayTennis": "Yes"},
					{"PlayTennis": "Yes"},
				},
				TargetColumn: "PlayTennis",
			},
			target:   "PlayTennis",
			expected: "Yes",
		},
		{
			name: "Two classes, majority Yes",
			dataset: &model.Dataset{
				RowInstances: []map[string]interface{}{
					{"PlayTennis": "Yes"},
					{"PlayTennis": "No"},
					{"PlayTennis": "Yes"},
				},
				TargetColumn: "PlayTennis",
			},
			target:   "PlayTennis",
			expected: "Yes",
		},
		{
			name: "Two classes, majority No",
			dataset: &model.Dataset{
				RowInstances: []map[string]interface{}{
					{"PlayTennis": "No"},
					{"PlayTennis": "Yes"},
					{"PlayTennis": "No"},
				},
				TargetColumn: "PlayTennis",
			},
			target:   "PlayTennis",
			expected: "No",
		},
		{
			name: "Empty dataset",
			dataset: &model.Dataset{
				RowInstances: []map[string]interface{}{},
				TargetColumn: "PlayTennis",
			},
			target:   "PlayTennis",
			expected: "",
		},
		{
			name: "Missing target column",
			dataset: &model.Dataset{
				RowInstances: []map[string]interface{}{
					{"Outlook": "Sunny"},
					{"Outlook": "Rain"},
				},
				TargetColumn: "PlayTennis", // Target column does not exist in records
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MajorityClass(tt.dataset, tt.target)
			if result != tt.expected {
				t.Errorf("MajorityClass() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
