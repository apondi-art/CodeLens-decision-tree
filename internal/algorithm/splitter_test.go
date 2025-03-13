package algorithm

import (
	"testing"
	"CodeLens-decision-tree/internal/model"
)

func TestSplitDataset(t *testing.T) {
	dataset := &model.Dataset{
		RowInstances: []map[string]interface{}{
			{"attribute": "value1"},
			{"attribute": "value2"},
			{"attribute": "value1"},
		},
		ColumnAttributes: map[string]*model.Attribute{
			"attribute": {Name: "attribute"},
		},
	}

	attr := &model.Attribute{Name: "attribute"}
	value := "value1"

	subset, err := SplitDataset(dataset, attr, value)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(subset.RowInstances) != 2 {
		t.Fatalf("expected 2 instances, got %d", len(subset.RowInstances))
	}
}

func TestSplitDatasetNilDataset(t *testing.T) {
	attr := &model.Attribute{Name: "attribute"}
	_, err := SplitDataset(nil, attr, "value")
	if err == nil {
		t.Fatal("expected error for nil dataset, got none")
	}
}

func TestSplitDatasetNilAttribute(t *testing.T) {
	dataset := &model.Dataset{
		RowInstances: []map[string]interface{}{
			{"attribute": "value1"},
		},
		ColumnAttributes: map[string]*model.Attribute{
			"attribute": {Name: "attribute"},
		},
	}
	_, err := SplitDataset(dataset, nil, "value")
	if err == nil {
		t.Fatal("expected error for nil attribute, got none")
	}
}

func TestFindBestNumericalSplit(t *testing.T) {
	dataset := &model.Dataset{
		RowInstances: []map[string]interface{}{
			{"numerical": 1.0, "target": "A"},
			{"numerical": 2.0, "target": "B"},
			{"numerical": 3.0, "target": "A"},
		},
		ColumnAttributes: map[string]*model.Attribute{
			"attribute": {Name: "attribute"},
		},
	}

	attr := &model.Attribute{Name: "numerical"}
	threshold, gain := FindBestNumericalSplit(dataset, attr, "target")
	if threshold == 0 || gain == 0 {
		t.Fatal("expected valid threshold and gain")
	}
}

func TestDistributeInstance(t *testing.T) {
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
	dataset := &model.Dataset{
		RowInstances: []map[string]interface{}{
			{"numerical": 1.0, "target": "A"},
			{"numerical": 2.0, "target": "B"},
			{"numerical": 3.0, "target": "A"},
		},
		ColumnAttributes: map[string]*model.Attribute{
			"attribute": {Name: "attribute"},
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
// test edge case or missing data
func TestSplitDatasetWithMissingAttribute(t *testing.T) {
	dataset := &model.Dataset{
		RowInstances: []map[string]interface{}{
			{"attribute": "A"},
			{"attribute": nil},
			{"attribute": "B"},
		},
	}
	subset, err := SplitDataset(dataset, &model.Attribute{Name: "attribute"}, "A")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(subset.RowInstances) != 1 {
		t.Errorf("Expected 1 instance, got %d", len(subset.RowInstances))
	}
}