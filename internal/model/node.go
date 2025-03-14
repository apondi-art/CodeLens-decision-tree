package model

// Node represents a node in the decision tree


// Node represents a single node in the decision tree.
type Node struct {
	Attribute      *Attribute      // The attribute used for splitting at this node
	SplitValue     interface{}     // The value of the attribute for this node
	Children       map[interface{}]*Node // Map of child nodes
	IsLeaf         bool            // Indicates if the node is a leaf
	PredictedClass string          // The predicted class for leaf nodes
	Class     string      // The class label for leaf nodes
	Left      *Node      // Pointer to the left child node
	Right     *Node      // Pointer to the right child node
	// Additional metadata for pruning and analysis:
	Depth             int            // Node depth in the tree
	SampleCount       int            // Number of training samples at this node
	ClassDistribution map[string]int  // Distribution of classes at this node
	ErrorEstimate     float64        // Estimated error for pruning
	GainRatio         float64        // Gain ratio that led to this split
}

// Predict classifies a single instance
func (n *Node) Predict(instance map[string]interface{}) string {
	return ""
}

// PredictWithMissingValues handles missing values during prediction
func (n *Node) PredictWithMissingValues(instance map[string]interface{}) map[string]float64 {
	return map[string]float64{}
}

// Consider adding:
// PrintSubtree returns a string representation of the subtree
func (n *Node) PrintSubtree(indent string) string {
	return ""
}
