package serialization

import (
	"encoding/json"
	"fmt"
	"os"

	"CodeLens-decision-tree/internal/model"
)

// SerializeTree serializes the DecisionTree and writes it to a JSON file.
func SerializeTree(root *model.Node, path string) error {
	treeJSON, err := NodeToJSON(root)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(treeJSON, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

// NodeToJSON recursively converts a Node into a JSON-compatible format.
func NodeToJSON(node *model.Node) (map[string]interface{}, error) {
	if node == nil {
		return nil, nil
	}

	// Create a map to represent the node
	nodeMap := map[string]interface{}{
		"splitValue":        node.SplitValue,
		"isLeaf":            node.IsLeaf,
		"predictedClass":    node.PredictedClass,
		"depth":             node.Depth,
		"sampleCount":       node.SampleCount,
		"classDistribution": node.ClassDistribution,
		"errorEstimate":     node.ErrorEstimate,
		"gainRatio":         node.GainRatio,
	}

	// Serialize attribute if it exists
	if node.Attribute != nil {
		nodeMap["attribute"] = node.Attribute.Name
	}

	// Serialize children if they exist
	if len(node.Children) > 0 {
		childrenMap := make(map[string]interface{})
		for key, child := range node.Children {
			childJSON, err := NodeToJSON(child)
			if err != nil {
				return nil, err
			}
			childrenMap[toString(key)] = childJSON
		}
		nodeMap["children"] = childrenMap
	}

	return nodeMap, nil
}

// toString converts an interface{} to string for JSON keys
func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int, float64, bool:
		return fmt.Sprintf("%v", v)
	default:
		return "unknown"
	}
}
