package model

// DecisionTree represents a C4.5 decision tree model
type DecisionTree struct {
	Root           *Node
	Attributes     []*Attribute
	TargetAttr     string
	AttributeTypes map[string]AttributeType
	// Consider adding these fields:
	Metadata        map[string]interface{} // For storing training data stats, feature importance, etc.
	MaxDepth        int                    // Maximum depth used during training
	MinSamplesSplit int                    // Minimum samples required for a split
	PruningEnabled  bool                   // Whether pruning was applied
}

// Predict classifies a single instance
func (t *DecisionTree) Predict(instance map[string]interface{}) string {
	if t.Root == nil {
		return ""
	}
	return t.Root.Predict(instance)
}

// BatchPredict classifies multiple instances
func (t *DecisionTree) BatchPredict(dataset *Dataset) []string {
	predictions := make([]string, len(dataset.RowInstances))

	for i, instance := range dataset.RowInstances {
		predictions[i] = t.Predict(instance)
	}

	return predictions
}

// PredictWithConfidence returns predicted class and confidence score
func (t *DecisionTree) PredictWithConfidence(instance map[string]interface{}) (string, float64) {
	if t.Root == nil {
		return "", 0
	}

	// Get class distribution
	probabilities := t.Root.PredictWithMissingValues(instance)

	// Find the class with the highest probability
	var bestClass string
	var bestProb float64

	for class, prob := range probabilities {
		if prob > bestProb {
			bestProb = prob
			bestClass = class
		}
	}

	return bestClass, bestProb
}

// GetFeatureImportance returns importance score for each attribute
func (t *DecisionTree) GetFeatureImportance() map[string]float64 {
	importance := make(map[string]float64)

	// Initialize with zero importance for all attributes
	for _, attr := range t.Attributes {
		importance[attr.Name] = 0
	}

	// Calculate importance based on gain ratio
	calculateImportance(t.Root, importance)

	// Normalize importance values
	total := 0.0
	for _, imp := range importance {
		total += imp
	}

	if total > 0 {
		for attr := range importance {
			importance[attr] /= total
		}
	}

	return importance
}

// Helper function to calculate feature importance
func calculateImportance(node *Node, importance map[string]float64) {
	if node == nil || node.IsLeaf {
		return
	}

	// Add the gain ratio of this node to the importance of its attribute
	importance[node.Attribute.Name] += node.GainRatio * float64(node.SampleCount)

	// Recursively calculate importance for children
	if node.Attribute.Type == Numeric {
		calculateImportance(node.Left, importance)
		calculateImportance(node.Right, importance)
	} else {
		for _, child := range node.Children {
			calculateImportance(child, importance)
		}
	}
}
