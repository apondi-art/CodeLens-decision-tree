package serialization

import (
	"encoding/json"
	"os"

	"CodeLens-decision-tree/internal/model"
)

// DeserializeTree loads a decision tree from a JSON file
func DeserializeTree(path string) (*model.DecisionTree, error) {
	// Read JSON file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}

	// Convert JSON to DecisionTree
	tree := &model.DecisionTree{}
	tree.Root, err = JSONToNode(jsonData["Root"].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	return tree, nil
}

// JSONToNode converts a JSON object to a DecisionTree Node
func JSONToNode(jsonData map[string]interface{}) (*model.Node, error) {
	node := &model.Node{
		Attribute: &model.Attribute{Name: jsonData["Attribute"].(string)},
		IsLeaf:    jsonData["IsLeaf"].(bool),
	}

	if predictedClass, ok := jsonData["PredictedClass"].(string); ok {
		node.PredictedClass = predictedClass
	}

	if splitValue, ok := jsonData["SplitValue"].(string); ok {
		node.SplitValue = splitValue
	}

	// Recursively populate child nodes
	childrenMap := make(map[string]*model.Node)
	if children, ok := jsonData["Children"].(map[string]interface{}); ok {
		for key, childData := range children {
			childDataMap, ok := childData.(map[string]interface{})
			if !ok {
			}

			childNode, err2 := JSONToNode(childDataMap)
			if err2 != nil {
				return nil, err2
			}
			childrenMap[key] = childNode
		}
	}

	node.Children = make(map[interface{}]*model.Node)
	for key, child := range childrenMap {
		node.Children[key] = child
	}

	return node, nil
}
