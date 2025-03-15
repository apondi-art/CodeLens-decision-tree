package serialization

import (
	"encoding/json"
	"fmt"
	"os"

	"CodeLens-decision-tree/internal/model"
)

// CustomNode is a temporary structure that matches the JSON format
type CustomNode struct {
	Attribute         string                 `json:"attribute"`
	SplitValue        interface{}            `json:"splitValue"`
	Children          map[string]*CustomNode `json:"children"`
	IsLeaf            bool                   `json:"isLeaf"`
	PredictedClass    string                 `json:"predictedClass"`
	ClassDistribution map[string]int         `json:"classDistribution"`
	Depth             int                    `json:"depth"`
	ErrorEstimate     float64                `json:"errorEstimate"`
	GainRatio         float64                `json:"gainRatio"`
	SampleCount       int                    `json:"sampleCount"`
}

// DeserializeTree loads a decision tree from a JSON file
func DeserializeTree(path string) (*model.DecisionTree, error) {
	// Read JSON file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read model file: %v", err)
	}

	// Parse JSON into our custom structure
	var customNode CustomNode
	if err := json.Unmarshal(data, &customNode); err != nil {
		return nil, fmt.Errorf("failed to parse model JSON: %v", err)
	}

	// Convert the custom node to a model.Node
	rootNode := convertToModelNode(&customNode)

	// Create tree with converted root node
	tree := &model.DecisionTree{
		Root:       rootNode,
		TargetAttr: customNode.Attribute, // Set the target attribute
	}

	return tree, nil
}

// Convert the CustomNode to a model.Node
func convertToModelNode(customNode *CustomNode) *model.Node {
	if customNode == nil {
		return nil
	}

	// Create an attribute structure
	attribute := &model.Attribute{
		Name: customNode.Attribute,
		// Assume it's a categorical attribute
		Type: model.Categorical,
	}

	// Create the node with proper attribute
	node := &model.Node{
		Attribute:         attribute,
		SplitValue:        customNode.SplitValue,
		IsLeaf:            customNode.IsLeaf,
		PredictedClass:    customNode.PredictedClass,
		Depth:             customNode.Depth,
		SampleCount:       customNode.SampleCount,
		ClassDistribution: customNode.ClassDistribution,
		ErrorEstimate:     customNode.ErrorEstimate,
		GainRatio:         customNode.GainRatio,
		Children:          make(map[interface{}]*model.Node),
	}

	// Convert children if any
	for key, childCustomNode := range customNode.Children {
		modelChild := convertToModelNode(childCustomNode)
		node.Children[key] = modelChild
	}

	return node
}
