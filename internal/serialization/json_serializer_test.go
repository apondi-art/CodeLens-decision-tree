package serialization

import (
	"encoding/json"
	"os"
	"testing"

	"CodeLens-decision-tree/internal/model"
)

// Helper function to create a sample decision tree
func createSampleTree() *model.Node {
	root := &model.Node{
		Attribute:         &model.Attribute{Name: "Color"},
		SplitValue:        "Red",
		IsLeaf:            false,
		Depth:             0,
		SampleCount:       100,
		ClassDistribution: map[string]int{"Yes": 60, "No": 40},
		ErrorEstimate:     0.2,
		GainRatio:         0.5,
		Children:          make(map[interface{}]*model.Node),
	}

	// Adding child nodes
	root.Children["Red"] = &model.Node{
		PredictedClass:    "Yes",
		IsLeaf:            true,
		Depth:             1,
		SampleCount:       60,
		ClassDistribution: map[string]int{"Yes": 60},
	}

	root.Children["Blue"] = &model.Node{
		PredictedClass:    "No",
		IsLeaf:            true,
		Depth:             1,
		SampleCount:       40,
		ClassDistribution: map[string]int{"No": 40},
	}

	return root
}

// Test NodeToJSON function
func TestNodeToJSON(t *testing.T) {
	tree := createSampleTree()

	jsonData, err := NodeToJSON(tree)
	if err != nil {
		t.Fatalf("Error serializing node: %v", err)
	}

	// Convert to actual JSON string
	data, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		t.Fatalf("Error marshalling JSON: %v", err)
	}

	t.Logf("Serialized JSON Output:\n%s", string(data))
}

// Test SerializeTree function
func TestSerializeTree(t *testing.T) {
	tree := createSampleTree()
	filePath := "test_tree.json"

	// Call serialization function
	err := SerializeTree(tree, filePath)
	if err != nil {
		t.Fatalf("Error writing JSON file: %v", err)
	}

	// Check if file exists
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Fatalf("Serialized file not found: %s", filePath)
	}

	// Clean up
	_ = os.Remove(filePath)
}
