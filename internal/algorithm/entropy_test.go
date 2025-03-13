package algorithm

import (
	"CodeLens-decision-tree/internal/model"
	"math"
	"testing"
)

// Helper function to compare floating-point numbers with a small tolerance.
func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

// TestCalculateEntropy verifies entropy calculations.
func TestCalculateEntropy(t *testing.T) {
	tests := []struct {
		name       string
		dataset    *model.Dataset
		targetAttr string
		expected   float64
	}{
		{
			name:       "Empty dataset",
			dataset:    &model.Dataset{Instances: []map[string]interface{}{}},
			targetAttr: "class",
			expected:   0.0,
		},
		{
			name: "Single class (pure dataset)",
			dataset: &model.Dataset{Instances: []map[string]interface{}{
				{"class": "A"},
				{"class": "A"},
				{"class": "A"},
			}},
			targetAttr: "class",
			expected:   0.0,
		},
		{
			name: "Two-class balanced dataset",
			dataset: &model.Dataset{Instances: []map[string]interface{}{
				{"class": "A"},
				{"class": "B"},
			}},
			targetAttr: "class",
			expected:   1.0,
		},
		{
			name: "Two-class imbalanced dataset",
			dataset: &model.Dataset{Instances: []map[string]interface{}{
				{"class": "A"},
				{"class": "A"},
				{"class": "B"},
				{"class": "B"},
				{"class": "B"},
			}},
			targetAttr: "class",
			expected:   0.97095, // Calculated manually
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CalculateEntropy(test.dataset, test.targetAttr)
			if !almostEqual(got, test.expected, 1e-4) {
				t.Errorf("Expected %f, got %f", test.expected, got)
			}
		})
	}
}

// TestCalculateInformationGain verifies information gain calculations.
func TestCalculateInformationGain(t *testing.T) {
	tests := []struct {
		name       string
		dataset    *model.Dataset
		attr       *model.Attribute
		targetAttr string
		expected   float64
	}{
		{
			name: "Zero information gain (no effect of split)",
			dataset: &model.Dataset{Instances: []map[string]interface{}{
				{"attr": "X", "class": "A"},
				{"attr": "X", "class": "A"},
			}},
			attr:       &model.Attribute{Name: "attr"},
			targetAttr: "class",
			expected:   0.0,
		},
		{
			name: "Full information gain (perfect split)",
			dataset: &model.Dataset{Instances: []map[string]interface{}{
				{"attr": "X", "class": "A"},
				{"attr": "Y", "class": "B"},
			}},
			attr:       &model.Attribute{Name: "attr"},
			targetAttr: "class",
			expected:   1.0,
		},
		{
			name: "Partial information gain",
			dataset: &model.Dataset{Instances: []map[string]interface{}{
				{"attr": "X", "class": "A"},
				{"attr": "X", "class": "A"},
				{"attr": "Y", "class": "B"},
				{"attr": "Y", "class": "B"},
				{"attr": "Y", "class": "A"},
			}},
			attr:       &model.Attribute{Name: "attr"},
			targetAttr: "class",
			expected:   0.17095, // Manually calculated
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CalculateInformationGain(test.dataset, test.attr, test.targetAttr)
			if !almostEqual(got, test.expected, 1e-4) {
				t.Errorf("Expected %f, got %f", test.expected, got)
			}
		})
	}
}

// TestCalculateGainRatio verifies gain ratio calculations.
func TestCalculateGainRatio(t *testing.T) {
	tests := []struct {
		name       string
		dataset    *model.Dataset
		attr       *model.Attribute
		targetAttr string
		expected   float64
	}{
		{
			name: "Zero gain ratio (no split effect)",
			dataset: &model.Dataset{Instances: []map[string]interface{}{
				{"attr": "X", "class": "A"},
				{"attr": "X", "class": "A"},
			}},
			attr:       &model.Attribute{Name: "attr"},
			targetAttr: "class",
			expected:   0.0,
		},
		{
			name: "High gain ratio (perfect split)",
			dataset: &model.Dataset{Instances: []map[string]interface{}{
				{"attr": "X", "class": "A"},
				{"attr": "Y", "class": "B"},
			}},
			attr:       &model.Attribute{Name: "attr"},
			targetAttr: "class",
			expected:   1.0,
		},
		{
			name: "Intermediate gain ratio",
			dataset: &model.Dataset{Instances: []map[string]interface{}{
				{"attr": "X", "class": "A"},
				{"attr": "X", "class": "A"},
				{"attr": "Y", "class": "B"},
				{"attr": "Y", "class": "B"},
				{"attr": "Y", "class": "A"},
			}},
			attr:       &model.Attribute{Name: "attr"},
			targetAttr: "class",
			expected:   0.2732, // Manually calculated
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CalculateGainRatio(test.dataset, test.attr, test.targetAttr)
			if !almostEqual(got, test.expected, 1e-4) {
				t.Errorf("Expected %f, got %f", test.expected, got)
			}
		})
	}
}
