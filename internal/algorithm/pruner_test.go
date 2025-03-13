package algorithm

import (
	"testing"
	"CodeLens-decision-tree/internal/model"
)

// TestPruneTree tests the PruneTree function.
func TestPruneTree(t *testing.T) {
	tests := []struct {
		name         string
		root         *model.Node
		validationSet *model.Dataset
		expectedClass string
	}{
		{
			name: "Nil root",
			root: nil,
			validationSet: &model.Dataset{RowInstances: []map[string]interface{}{}},
			expectedClass: "", // Expect nil or no class
		},
		{
			name: "Single leaf node",
			root: &model.Node{Class: "A"},
			validationSet: &model.Dataset{RowInstances: []map[string]interface{}{}},
			expectedClass: "A", // Expect the same class
		},
		{
			name: "Pruning with majority class",
			root: &model.Node{
				Class: "B",
				Left: &model.Node{Class: "A"},
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
				Left: &model.Node{Class: "B"},
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