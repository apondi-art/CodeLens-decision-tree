package serialization

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeserializeTree(t *testing.T) {
	// Create a temporary JSON file for testing
	tempFile, err := os.CreateTemp("", "test_tree_*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write sample JSON data to the file
	jsonData := `{
		"Root": {
			"Attribute": "Color",
			"IsLeaf": false,
			"Children": {
				"Red": {
					"Attribute": "Size",
					"IsLeaf": true,
					"PredictedClass": "Apple"
				}
			}
		}
	}`
	_, err = tempFile.WriteString(jsonData)
	assert.NoError(t, err)
	tempFile.Close()

	// Call DeserializeTree
	tree, err := DeserializeTree(tempFile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, tree)

	// Validate the tree structure
	assert.Equal(t, "Color", tree.Root.Attribute.Name)
	assert.False(t, tree.Root.IsLeaf)
	assert.Contains(t, tree.Root.Children, "Red")

	redNode := tree.Root.Children["Red"]
	assert.Equal(t, "Size", redNode.Attribute.Name)
	assert.True(t, redNode.IsLeaf)
	assert.Equal(t, "Apple", redNode.PredictedClass)
}

func TestJSONToNode(t *testing.T) {
	// Sample JSON data for a node
	jsonData := map[string]interface{}{
		"Attribute": "Color",
		"IsLeaf":    false,
		"Children": map[string]interface{}{
			"Red": map[string]interface{}{
				"Attribute":      "Size",
				"IsLeaf":         true,
				"PredictedClass": "Apple",
			},
		},
	}

	// Call JSONToNode
	node, err := JSONToNode(jsonData)
	assert.NoError(t, err)
	assert.NotNil(t, node)

	// Validate the node
	assert.Equal(t, "Color", node.Attribute.Name)
	assert.False(t, node.IsLeaf)

	// Validate child nodes
	assert.Contains(t, node.Children, "Red")

	redNode := node.Children["Red"]
	assert.Equal(t, "Size", redNode.Attribute.Name)
	assert.True(t, redNode.IsLeaf)
	assert.Equal(t, "Apple", redNode.PredictedClass)
}
