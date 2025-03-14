package algorithm

import (
	"testing"

	"CodeLens-decision-tree/internal/model"
)

func TestBuildTree(t *testing.T) {
	tests := []struct {
		name          string
		dataset       *model.Dataset
		attributes    []*model.Attribute
		targetAttr    string
		maxDepth      int
		expectedClass string // Expected class for a sample instance
		expectedError bool   // Whether an error is expected
	}{
		{
			name: "Pure dataset",
			dataset: &model.Dataset{
				RowInstances: []map[string]interface{}{
					{"Outlook": "Sunny", "PlayTennis": "No"},
					{"Outlook": "Sunny", "PlayTennis": "No"},
				},
				TargetColumn: "PlayTennis",
			},
			attributes: []*model.Attribute{
				{Name: "Outlook", Type: model.Categorical},
			},
			targetAttr:    "PlayTennis",
			maxDepth:      5,
			expectedClass: "No", // All instances belong to the same class
			expectedError: false,
		},
		{
			name: "Numerical attribute",
			dataset: &model.Dataset{
				RowInstances: []map[string]interface{}{
					{"Temperature": 85.0, "PlayTennis": "No"},
					{"Temperature": 80.0, "PlayTennis": "Yes"},
					{"Temperature": 83.0, "PlayTennis": "No"},
				},
				TargetColumn: "PlayTennis",
			},
			attributes: []*model.Attribute{
				{Name: "Temperature", Type: model.Numeric},
			},
			targetAttr:    "PlayTennis",
			maxDepth:      5,
			expectedClass: "No", // Temperature <= 82.5 -> No
			expectedError: false,
		},
		{
			name: "Maximum depth reached",
			dataset: &model.Dataset{
				RowInstances: []map[string]interface{}{
					{"Outlook": "Sunny", "PlayTennis": "No"},
					{"Outlook": "Rain", "PlayTennis": "Yes"},
					{"Outlook": "Sunny", "PlayTennis": "Yes"},
					{"Outlook": "Rain", "PlayTennis": "No"},
					{"Outlook": "Overcast", "PlayTennis": "Yes"},
					{"Outlook": "Sunny", "PlayTennis": "Yes"},
					{"Outlook": "Overcast", "PlayTennis": "No"},
					{"Outlook": "Rain", "PlayTennis": "Yes"},
				},
				TargetColumn: "PlayTennis",
			},
			attributes: []*model.Attribute{
				{Name: "Outlook", Type: model.Categorical},
			},
			targetAttr:    "PlayTennis",
			maxDepth:      0,     // Tree stops at depth 0
			expectedClass: "Yes", // Majority class at root
			expectedError: false,
		},
		{
			name: "Empty dataset",
			dataset: &model.Dataset{
				RowInstances: []map[string]interface{}{},
				TargetColumn: "PlayTennis",
			},
			attributes: []*model.Attribute{
				{Name: "Outlook", Type: model.Categorical},
			},
			targetAttr:    "PlayTennis",
			maxDepth:      5,
			expectedClass: "", // No instances, so no prediction
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := BuildTree(tt.dataset, tt.attributes, tt.targetAttr, 0, tt.maxDepth)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("BuildTree() error = %v, expected no error", err)
				return
			}

			// Test prediction for a sample instance.
			var instance map[string]interface{}
			if len(tt.dataset.RowInstances) > 0 {
				instance = tt.dataset.RowInstances[0]
			} else {
				instance = map[string]interface{}{}
			}

			predictedClass := tree.Predict(instance)
			if predictedClass != tt.expectedClass {
				t.Errorf("Predicted class = %v, expected %v", predictedClass, tt.expectedClass)
			}
		})
	}
}

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
