package serialization

import (
	"os"
	"testing"

	"CodeLens-decision-tree/internal/model"
)

// Test SerializeTree function
func TestSerializeTree(t *testing.T) {
	// Helper function to create a sample decision tree
	createSampleTree := func() *model.TreeNode {
		root := &model.TreeNode{
			SplitAttr:         "Color",
			SplitValue:        "Red",
			IsLeaf:            false,
			ClassDistribution: map[string]float64{"Yes": 60, "No": 40},
			Children:          make(map[string]*model.TreeNode),
		}

		// Adding child nodes
		root.Children["Red"] = &model.TreeNode{
			PredictedClass:    "Yes",
			IsLeaf:            true,
			ClassDistribution: map[string]float64{"Yes": 60},
		}

		root.Children["Blue"] = &model.TreeNode{
			PredictedClass:    "No",
			IsLeaf:            true,
			ClassDistribution: map[string]float64{"No": 40},
		}

		return root
	}
	// Create a sample decision tree
	dt := &model.DecisionTree{
		Root:       createSampleTree(),
		TargetAttr: "Species",
	}
	filePath := "test_tree.json"

	// Call serialization function
	err := SerializeTree(dt, filePath)
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
