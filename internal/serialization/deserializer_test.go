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
		"attribute": "Color",
		"isLeaf": false,
		"children": {
			"Red": {
				"attribute": "Size",
				"isLeaf": true,
				"predictedClass": "Apple"
			}
		},
		"classDistribution": {"Apple": 10, "Banana": 5},
		"depth": 1,
		"errorEstimate": 0.1,
		"gainRatio": 0.25,
		"sampleCount": 15
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

func TestConvertToModelNode(t *testing.T) {
	// Sample CustomNode data
	customNode := &CustomNode{
		Attribute:      "Color",
		SplitValue:     "Red",
		IsLeaf:         false,
		PredictedClass: "",
		ClassDistribution: map[string]int{
			"Apple":  10,
			"Banana": 5,
		},
		Depth:         1,
		ErrorEstimate: 0.1,
		GainRatio:     0.25,
		SampleCount:   15,
		Children: map[string]*CustomNode{
			"Red": {
				Attribute:      "Size",
				IsLeaf:         true,
				PredictedClass: "Apple",
			},
		},
	}

	// Call convertToModelNode
	node := convertToModelNode(customNode)

	// Validate the node
	assert.Equal(t, "Color", node.Attribute.Name)
	assert.False(t, node.IsLeaf)
	assert.Equal(t, "Red", node.SplitValue)
	assert.Equal(t, 1, node.Depth)
	assert.Equal(t, 0.1, node.ErrorEstimate)
	assert.Equal(t, 0.25, node.GainRatio)
	assert.Equal(t, 15, node.SampleCount)

	// Validate child nodes
	assert.Contains(t, node.Children, "Red")

	redNode := node.Children["Red"]
	assert.Equal(t, "Size", redNode.Attribute.Name)
	assert.True(t, redNode.IsLeaf)
	assert.Equal(t, "Apple", redNode.PredictedClass)
}
