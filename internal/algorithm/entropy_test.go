package algorithm

import (
	"math"
	"testing"

	"CodeLens-decision-tree/internal/model"
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
			dataset:    &model.Dataset{RowInstances: []map[string]interface{}{}},
			targetAttr: "class",
			expected:   0.0,
		},
		{
			name: "Single class (pure dataset)",
			dataset: &model.Dataset{RowInstances: []map[string]interface{}{
				{"class": "A"},
				{"class": "A"},
				{"class": "A"},
			}},
			targetAttr: "class",
			expected:   0.0,
		},
		{
			name: "Two-class balanced dataset",
			dataset: &model.Dataset{RowInstances: []map[string]interface{}{
				{"class": "A"},
				{"class": "B"},
			}},
			targetAttr: "class",
			expected:   1.0,
		},
		{
			name: "Two-class imbalanced dataset",
			dataset: &model.Dataset{RowInstances: []map[string]interface{}{
				{"class": "A"},
				{"class": "A"},
				{"class": "B"},
				{"class": "B"},
				{"class": "B"},
			}},
			targetAttr: "class",
			expected:   0.97095, // Calculated manually
		},
		{
			name: "Missing values in target attribute",
			dataset: &model.Dataset{
				RowInstances: []map[string]interface{}{
					{"class": "A"},
					{"class": nil}, // Missing value
					{"class": "B"},
				},
			},
			targetAttr: "class",
			expected:   1.0, // Update based on your calculation logic
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
			dataset: &model.Dataset{RowInstances: []map[string]interface{}{
				{"attr": "X", "class": "A"},
				{"attr": "X", "class": "A"},
			}},
			attr:       &model.Attribute{Name: "attr"},
			targetAttr: "class",
			expected:   0.0,
		},
		{
			name: "Full information gain (perfect split)",
			dataset: &model.Dataset{RowInstances: []map[string]interface{}{
				{"attr": "X", "class": "A"},
				{"attr": "Y", "class": "B"},
			}},
			attr:       &model.Attribute{Name: "attr"},
			targetAttr: "class",
			expected:   1.0,
		},
		{
			name: "Partial information gain",
			dataset: &model.Dataset{RowInstances: []map[string]interface{}{
				{"attr": "X", "class": "A"},
				{"attr": "X", "class": "A"},
				{"attr": "Y", "class": "B"},
				{"attr": "Y", "class": "B"},
				{"attr": "Y", "class": "A"},
			}},
			attr:       &model.Attribute{Name: "attr"},
			targetAttr: "class",
			expected:   0.419973, // Updated expected value
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
			dataset: &model.Dataset{RowInstances: []map[string]interface{}{
				{"attr": "X", "class": "A"},
				{"attr": "X", "class": "A"},
			}},
			attr:       &model.Attribute{Name: "attr"},
			targetAttr: "class",
			expected:   0.0,
		},
		{
			name: "High gain ratio (perfect split)",
			dataset: &model.Dataset{RowInstances: []map[string]interface{}{
				{"attr": "X", "class": "A"},
				{"attr": "Y", "class": "B"},
			}},
			attr:       &model.Attribute{Name: "attr"},
			targetAttr: "class",
			expected:   1.0,
		},
		{
			name: "Intermediate gain ratio",
			dataset: &model.Dataset{RowInstances: []map[string]interface{}{
				{"attr": "X", "class": "A"},
				{"attr": "X", "class": "A"},
				{"attr": "Y", "class": "B"},
				{"attr": "Y", "class": "B"},
				{"attr": "Y", "class": "A"},
			}},
			attr:       &model.Attribute{Name: "attr"},
			targetAttr: "class",
			expected:   0.432538, // Updated expected value
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

// TestCalculateEntropyMissingValues checks how the function handles missing values.
func TestCalculateEntropyWithMissingValues(t *testing.T) {
	dataset := &model.Dataset{
		RowInstances: []map[string]interface{}{
			{"attribute1": "A", "target": "yes"},
			{"attribute1": nil, "target": "no"},
			{"attribute1": "B", "target": "yes"},
		},
	}

	// Calculate expected entropy
	expectedEntropy := -(2.0/3.0*math.Log2(2.0/3.0) + 1.0/3.0*math.Log2(1.0/3.0))

	actualEntropy := CalculateEntropy(dataset, "target")
	if actualEntropy != expectedEntropy {
		t.Errorf("Expected %v, got %v", expectedEntropy, actualEntropy)
	}
}
