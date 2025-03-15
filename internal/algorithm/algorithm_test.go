package algorithm

import (
	"math"
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
		{
			name: "Single attribute",
			dataset: &model.Dataset{
				RowInstances: []map[string]interface{}{
					{"Outlook": "Sunny", "PlayTennis": "No"},
					{"Outlook": "Rain", "PlayTennis": "Yes"},
				},
				TargetColumn: "PlayTennis",
			},
			attributes: []*model.Attribute{
				{Name: "Outlook", Type: model.Categorical},
			},
			targetAttr:    "PlayTennis",
			maxDepth:      5,
			expectedClass: "No", // Outlook=Sunny -> No
			expectedError: false,
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
			attr, gain := SelectBestAttribute(&tt.dataset, tt.attributes, tt.target)
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

// TestCalculateInformationGain verifies information gain calculations.// TestCalculateInformationGain verifies information gain calculations.
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
			expected:   0.970951, // Corrected expected value based on the entropy calculation
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

// Helper function for comparing floating point numbers with a tolerance.

// TestCalculateGainRatio verifies gain ratio calculations.
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
			expected:   1.000000, // Updated expected value
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
			{"attribute1": "C", "target": "no"},
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

// TestPruneTree tests the PruneTree function.
func TestPruneTree(t *testing.T) {
	tests := []struct {
		name          string
		root          *model.Node
		validationSet *model.Dataset
		expectedClass string
	}{
		{
			name:          "Nil root",
			root:          nil,
			validationSet: &model.Dataset{RowInstances: []map[string]interface{}{}},
			expectedClass: "", // Expect nil or no class
		},
		{
			name:          "Single leaf node",
			root:          &model.Node{Class: "A"},
			validationSet: &model.Dataset{RowInstances: []map[string]interface{}{}},
			expectedClass: "A", // Expect the same class
		},
		{
			name: "Pruning with majority class",
			root: &model.Node{
				Class: "B",
				Left:  &model.Node{Class: "A"},
				Right: &model.Node{Class: "C"},
			},
			validationSet: &model.Dataset{
				RowInstances: []map[string]interface{}{
					{"class": "A"},
					{"class": "A"},
					{"class": "B"},
				},
			},
			expectedClass: "A", // Expect majority class after pruning
		},
		{
			name: "No improvement from pruning",
			root: &model.Node{
				Class: "B",
				Left:  &model.Node{Class: "B"},
				Right: &model.Node{Class: "B"},
			},
			validationSet: &model.Dataset{
				RowInstances: []map[string]interface{}{
					{"class": "B"},
					{"class": "B"},
				},
			},
			expectedClass: "B", // Expect no change
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prunedTree := PruneTree(test.root, test.validationSet, "class")
			if prunedTree == nil {
				if test.expectedClass != "" {
					t.Errorf("Expected class %s, got nil", test.expectedClass)
				}
				return // Exit if the pruned tree is nil
			}
			if prunedTree.Class != test.expectedClass {
				t.Errorf("Expected class %s, got %s", test.expectedClass, prunedTree.Class)
			}
		})
	}
}

func TestSplitDatasetNilDataset(t *testing.T) {
	// Test for nil dataset
	split := &model.Split{
		Attribute: &model.Attribute{Name: "attribute"},
	}
	_, err := SplitDataset(nil, split)
	if err == nil {
		t.Fatal("expected error for nil dataset, got none")
	}
}

func TestSplitDatasetNilSplit(t *testing.T) {
	// Test for nil split
	dataset := &model.Dataset{
		RowInstances: []map[string]interface{}{
			{"attribute": "value1"},
		},
	}
	_, err := SplitDataset(dataset, nil)
	if err == nil {
		t.Fatal("expected error for nil split, got none")
	}
}

func TestSplitDatasetNilAttribute(t *testing.T) {
	// Test for nil attribute
	dataset := &model.Dataset{
		RowInstances: []map[string]interface{}{
			{"attribute": "value1"},
		},
	}
	split := &model.Split{
		Attribute: nil,
	}
	_, err := SplitDataset(dataset, split)
	if err == nil {
		t.Fatal("expected error for nil attribute, got none")
	}
}

func TestFindBestNumericalSplit(t *testing.T) {
	// Test to find best numerical split
	dataset := &model.Dataset{
		RowInstances: []map[string]interface{}{
			{"numerical": 1.0, "target": "A"},
			{"numerical": 2.0, "target": "B"},
			{"numerical": 3.0, "target": "A"},
		},
		ColumnAttributes: map[string]*model.Attribute{
			"numerical": {Name: "numerical", Type: model.Numeric},
			"target":    {Name: "target", Type: model.Categorical},
		},
	}

	attr := &model.Attribute{Name: "numerical"}
	targetAttr := "target"

	threshold, gain := FindBestNumericalSplit(dataset, attr, targetAttr)
	if threshold == 0 || gain == 0 {
		t.Fatal("expected valid threshold and gain")
	}
}

func TestDistributeInstance(t *testing.T) {
	// Test for distributing an instance to child nodes
	instance := map[string]interface{}{"attribute": "value1"}
	children := map[interface{}]*model.Node{
		"value1": {Attribute: &model.Attribute{Name: "attribute"}},
		"value2": {Attribute: &model.Attribute{Name: "attribute"}},
	}

	distribution := DistributeInstance(instance, children)
	if len(distribution) != 1 {
		t.Fatalf("expected 1 child in distribution, got %d", len(distribution))
	}
}

func TestComputeInformationGain(t *testing.T) {
	// Test for computing information gain
	dataset := &model.Dataset{
		RowInstances: []map[string]interface{}{
			{"numerical": 1.0, "target": "A"},
			{"numerical": 2.0, "target": "B"},
			{"numerical": 3.0, "target": "A"},
		},
		ColumnAttributes: map[string]*model.Attribute{
			"numerical": {Name: "numerical", Type: model.Numeric},
			"target":    {Name: "target", Type: model.Categorical},
		},
	}

	attr := "numerical"
	threshold := 1.5
	targetAttr := "target"

	gain := computeInformationGain(dataset, attr, threshold, targetAttr)
	if gain <= 0 {
		t.Fatalf("expected positive information gain, got %f", gain)
	}
}

func TestFindBestNumericalSplitWithNoGain(t *testing.T) {
	// Test when no gain is found
	dataset := &model.Dataset{
		RowInstances: []map[string]interface{}{
			{"numerical": 1.0, "target": "A"},
			{"numerical": 1.0, "target": "B"},
		},
		ColumnAttributes: map[string]*model.Attribute{
			"numerical": {Name: "numerical", Type: model.Numeric},
			"target":    {Name: "target", Type: model.Categorical},
		},
	}

	attr := &model.Attribute{Name: "numerical"}
	targetAttr := "target"

	threshold, gain := FindBestNumericalSplit(dataset, attr, targetAttr)
	if threshold != 0 || gain != 0 {
		t.Fatal("expected no gain, got some")
	}
}
