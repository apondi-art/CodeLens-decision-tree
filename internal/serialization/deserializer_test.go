package serialization

import (
	"os"
	"testing"
)

func TestDeserializeTree(t *testing.T) {
	// Create a temporary JSON file for testing
	tempFile, err := os.CreateTemp("", "test_tree_*.json")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write sample JSON data to the file
	jsonData := `{
		"Root": {
			"SplitAttr": "Color",
			"SplitValue": null,
			"IsLeaf": false,
			"ClassDistribution": {
				"Yes": 60,
				"No": 40
			},
			"Children": {
				"0": {
					"SplitAttr": "Size",
					"IsLeaf": true,
					"PredictedClass": "Apple",
					"ClassDistribution": null,
					"SplitValue": null,
					"Children": null
				}
			}
		},
      "TargetAttr": "Target"
	}`
	_, err = tempFile.WriteString(jsonData)
	if err != nil {
		t.Fatalf("Error writing JSON data to file: %v", err)
	}
	tempFile.Close()

	// Call DeserializeTree
	tree, err := DeserializeTree(tempFile.Name())
	if err != nil {
		t.Fatalf("Error deserializing tree: %v", err)
	}
	if tree == nil {
		t.Fatalf("Deserialized tree is nil")
	}

	// Validate the tree structure
	if tree.Root.SplitAttr != "Color" {
		t.Errorf("Expected SplitAttr to be Color, but got %s", tree.Root.SplitAttr)
	}
	if tree.Root.IsLeaf {
		t.Errorf("Expected IsLeaf to be false, but got true")
	}
	if tree.TargetAttr != "Target" {
		t.Errorf("Expected TargetAttr to be Target, but got %s", tree.TargetAttr)
	}
	if _, ok := tree.Root.Children["0"]; !ok {
		t.Errorf("Expected Children to contain key 0")
	}

	redNode := tree.Root.Children["0"]
	if redNode.SplitAttr != "Size" {
		t.Errorf("Expected SplitAttr to be Size, but got %s", redNode.SplitAttr)
	}
	if redNode.IsLeaf != true {
		t.Errorf("Expected IsLeaf to be true, but got false")
	}
	if redNode.PredictedClass != "Apple" {
		t.Errorf("Expected PredictedClass to be Apple, but got %s", redNode.PredictedClass)
	}
}
