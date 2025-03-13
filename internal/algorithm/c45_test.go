package algorithm

import (
	"testing"

	"CodeLens-decision-tree/internal/model"
)

func TestSelectBestAttribute(t *testing.T) {
	tests := []struct {
		name         string
		dataset      model.Dataset
		attributes   []*model.Attribute
		target       string
		expectedAttr string
		expectedGain float64
	}{
		{
			name: "Best attribute is Outlook",
			dataset: model.Dataset{
				RowInstances: []map[string]interface{}{
					{"Outlook": "Sunny", "PlayTennis": "No"},
					{"Outlook": "Sunny", "PlayTennis": "No"},
					{"Outlook": "Overcast", "PlayTennis": "Yes"},
					{"Outlook": "Rain", "PlayTennis": "Yes"},
				},
				TargetColumn: "PlayTennis",
			},
			attributes: []*model.Attribute{
				{Name: "Outlook", Type: model.Categorical},
				{Name: "Temperature", Type: model.Categorical},
			},
			target:       "PlayTennis",
			expectedAttr: "Outlook",
			expectedGain: 0.666667, // Precomputed gain ratio for Outlook
		},
		{
			name: "No attributes left",
			dataset: model.Dataset{
				RowInstances: []map[string]interface{}{
					{"PlayTennis": "Yes"},
					{"PlayTennis": "Yes"},
				},
				TargetColumn: "PlayTennis",
			},
			attributes:   []*model.Attribute{},
			target:       "PlayTennis",
			expectedAttr: "",
			expectedGain: -1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr, gain := SelectBestAttribute(tt.dataset, tt.attributes, tt.target)
			if attr != nil && attr.Name != tt.expectedAttr {
				t.Errorf("SelectBestAttribute() attribute = %v, expected %v", attr.Name, tt.expectedAttr)
			}
			if gain != tt.expectedGain {
				t.Errorf("SelectBestAttribute() gain = %v, expected %v", gain, tt.expectedGain)
			}
		})
	}
}

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
		{
			name: "All missing values in target column",
			dataset: &model.Dataset{
				RowInstances: []map[string]interface{}{
					{"PlayTennis": nil},
					{"PlayTennis": nil},
				},
				TargetColumn: "PlayTennis",
			},
			target:   "PlayTennis",
			expected: "",
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
