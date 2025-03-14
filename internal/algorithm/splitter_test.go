package algorithm

import (
	"CodeLens-decision-tree/internal/model"
	"testing"
)

func TestSplitDataset(t *testing.T) {
	// Test for correct splitting by a given attribute and value
	dataset := &model.Dataset{
		RowInstances: []map[string]interface{}{
			{"attribute": "value1"},
			{"attribute": "value2"},
			{"attribute": "value1"},
		},
		ColumnAttributes: map[string]*model.Attribute{
			"attribute": {Name: "attribute", Type: model.Categorical},
		},
	}

	split := &model.Split{
		Attribute:    &model.Attribute{Name: "attribute"},
		CategoricalMap: map[string]bool{"value1": true},
	}

	subset, err := SplitDataset(dataset, split)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(subset.RowInstances) != 2 {
		t.Fatalf("expected 2 instances, got %d", len(subset.RowInstances))
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

func TestSplitDatasetWithMissingAttribute(t *testing.T) {
	// Test for dataset with missing attribute values
	dataset := &model.Dataset{
		RowInstances: []map[string]interface{}{
			{"attribute": "A"},
			{"attribute": nil},
			{"attribute": "B"},
		},
		ColumnAttributes: map[string]*model.Attribute{
			"attribute": {Name: "attribute", Type: model.Categorical},
		},
	}

	split := &model.Split{
		Attribute:    &model.Attribute{Name: "attribute"},
		CategoricalMap: map[string]bool{"A": true},
	}

	subset, err := SplitDataset(dataset, split)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(subset.RowInstances) != 1 {
		t.Errorf("Expected 1 instance, got %d", len(subset.RowInstances))
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
