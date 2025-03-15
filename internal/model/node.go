package model

import (
	"fmt"
	"strconv"
	"strings"
)

// Node represents a single node in the decision tree.
type Node struct {
	Attribute      *Attribute            // The attribute used for splitting at this node
	SplitValue     interface{}           // The value of the attribute for this node
	Children       map[interface{}]*Node // Map of child nodes
	IsLeaf         bool                  // Indicates if the node is a leaf
	PredictedClass string                // The predicted class for leaf nodes
	Class          string                // The class label for leaf nodes
	Left           *Node                 // Pointer to the left child node
	Right          *Node                 // Pointer to the right child node
	// Additional metadata for pruning and analysis:
	Depth             int            // Node depth in the tree
	SampleCount       int            // Number of training samples at this node
	ClassDistribution map[string]int // Distribution of classes at this node
	ErrorEstimate     float64        // Estimated error for pruning
	GainRatio         float64        // Gain ratio that led to this split
}

// Predict classifies a single instance
func (n *Node) Predict(instance map[string]interface{}) string {
	// If this is a leaf node, return the predicted class
	if n.IsLeaf {
		return n.PredictedClass
	}

	// Get the attribute value from the instance
	value, exists := instance[n.Attribute.Name]

	// If the attribute doesn't exist in the instance or is nil
	if !exists || value == nil {
		// Handle missing values by using the majority class at this node
		return n.PredictedClass
	}

	// For numerical attributes, handle the comparison
	if n.Attribute.Type == Numeric {
		// Convert values to float64 for comparison
		instanceValue, ok := toFloat64(value)
		if !ok {
			// If conversion fails, return the majority class
			return n.PredictedClass
		}

		splitValue, ok := toFloat64(n.SplitValue)
		if !ok {
			// If conversion fails, return the majority class
			return n.PredictedClass
		}

		// Navigate to the appropriate child based on the comparison
		if instanceValue <= splitValue {
			if n.Left != nil {
				return n.Left.Predict(instance)
			}
		} else {
			if n.Right != nil {
				return n.Right.Predict(instance)
			}
		}

		// If the appropriate child doesn't exist, return the majority class
		return n.PredictedClass
	} else {
		// For categorical attributes, find the matching child
		childNode, exists := n.Children[value]
		if !exists {
			// If no matching child, return the majority class
			return n.PredictedClass
		}
		return childNode.Predict(instance)
	}
}

// Helper function to convert interface{} to float64
func toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case string:
		// Try to parse the string as a float
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	default:
		return 0, false
	}
}

// PredictWithMissingValues handles missing values during prediction by exploring all paths
// and returning a probability distribution over classes
func (n *Node) PredictWithMissingValues(instance map[string]interface{}) map[string]float64 {
	// If this is a leaf node, return 100% probability for the predicted class
	if n.IsLeaf {
		result := make(map[string]float64)
		result[n.PredictedClass] = 1.0
		return result
	}

	// Check if the attribute value is missing
	value, exists := instance[n.Attribute.Name]
	if !exists || value == nil {
		// If missing, explore all possible paths and merge results
		result := make(map[string]float64)

		// Calculate weights based on the number of samples in each branch
		totalSamples := 0
		for _, child := range n.Children {
			totalSamples += child.SampleCount
		}

		// Explore each child and merge results with appropriate weights
		for _, child := range n.Children {
			weight := float64(child.SampleCount) / float64(totalSamples)
			childPredictions := child.PredictWithMissingValues(instance)

			// Merge predictions with weights
			for class, prob := range childPredictions {
				result[class] += prob * weight
			}
		}

		return result
	}

	// For numerical attributes with existing values
	if n.Attribute.Type == Numeric {
		instanceValue, ok := toFloat64(value)
		if !ok {
			// If conversion fails, treat as missing
			return n.PredictWithMissingValues(map[string]interface{}{})
		}

		splitValue, ok := toFloat64(n.SplitValue)
		if !ok {
			// If conversion fails, use majority class
			result := make(map[string]float64)
			result[n.PredictedClass] = 1.0
			return result
		}

		// Navigate to appropriate child
		if instanceValue <= splitValue && n.Left != nil {
			return n.Left.PredictWithMissingValues(instance)
		} else if instanceValue > splitValue && n.Right != nil {
			return n.Right.PredictWithMissingValues(instance)
		}

		// If no appropriate child, return majority class
		result := make(map[string]float64)
		result[n.PredictedClass] = 1.0
		return result
	} else {
		// For categorical attributes with existing values
		childNode, exists := n.Children[value]
		if !exists {
			// If no matching child, return majority class
			result := make(map[string]float64)
			result[n.PredictedClass] = 1.0
			return result
		}
		return childNode.PredictWithMissingValues(instance)
	}
}

// PrintSubtree returns a string representation of the subtree
func (n *Node) PrintSubtree(indent string) string {
	if n == nil {
		return indent + "nil"
	}

	var sb strings.Builder

	if n.IsLeaf {
		sb.WriteString(fmt.Sprintf("%sLeaf: %s (Samples: %d)\n", indent, n.PredictedClass, n.SampleCount))
	} else {
		if n.Attribute != nil {
			if n.Attribute.Type == Numeric {
				sb.WriteString(fmt.Sprintf("%sNode: %s <= %v (Gain: %.4f, Samples: %d)\n",
					indent, n.Attribute.Name, n.SplitValue, n.GainRatio, n.SampleCount))
				sb.WriteString(fmt.Sprintf("%s  Left: %s\n", indent, n.Left.PrintSubtree(indent+"  ")))
				sb.WriteString(fmt.Sprintf("%s  Right: %s\n", indent, n.Right.PrintSubtree(indent+"  ")))
			} else {
				sb.WriteString(fmt.Sprintf("%sNode: %s (Gain: %.4f, Samples: %d)\n",
					indent, n.Attribute.Name, n.GainRatio, n.SampleCount))
				for value, child := range n.Children {
					sb.WriteString(fmt.Sprintf("%s  Value: %v -> %s\n", indent, value, child.PrintSubtree(indent+"  ")))
				}
			}
		} else {
			sb.WriteString(fmt.Sprintf("%sNode: (unknown attribute)\n", indent))
		}
	}

	return sb.String()
}
