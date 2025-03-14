package serialization

import (
	"CodeLens-decision-tree/internal/model"
	"encoding/json"
	"fmt"
	"os"
)

// DeserializeTree loads a decision tree from a JSON file
func DeserializeTree(path string) (*model.DecisionTree, error) {
	// Read JSON file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var tree model.DecisionTree
	type Alias model.DecisionTree
	err = json.Unmarshal(data, &struct {
		*Alias
	}{
		Alias: (*Alias)(&tree),
	})
	if err != nil {
		return nil, err
	}

	return &tree, nil
}

// JSONToTreeNode converts a JSON object to a DecisionTree Node
func JSONToTreeNode(jsonData map[string]interface{}) (*model.TreeNode, error) {
	node := &model.TreeNode{}

	if splitAttr, ok := jsonData["SplitAttr"].(string); ok {
		node.SplitAttr = splitAttr
	}

	if isLeaf, ok := jsonData["IsLeaf"].(bool); ok {
		node.IsLeaf = isLeaf
	}

	if predictedClass, ok := jsonData["PredictedClass"].(string); ok {
		node.PredictedClass = predictedClass
	}

	if splitValue, ok := jsonData["SplitValue"].(interface{}); ok {
		node.SplitValue = splitValue
	}

	if classDistribution, ok := jsonData["ClassDistribution"].(map[string]interface{}); ok {
		node.ClassDistribution = make(map[string]float64)
		for k, v := range classDistribution {
			if floatValue, ok := v.(float64); ok {
				node.ClassDistribution[k] = floatValue
			}
		}
	}

	// Recursively populate child nodes
	node.Children = make(map[string]*model.TreeNode)
	if children, ok := jsonData["Children"].([]interface{}); ok {
		for i, childData := range children {
			if childDataMap, ok := childData.(map[string]interface{}); ok {
				childNode, err := JSONToTreeNode(childDataMap)
				if err != nil {
					return nil, fmt.Errorf("error deserializing child %d: %w", i, err)
				}
				node.Children[fmt.Sprintf("%d", i)] = childNode
			}
		}
	}

	return node, nil
}
