package algorithm

import (
	"CodeLens-decision-tree/internal/model"
)

// PruneTree prunes the decision tree based on the validation set.
// It recursively traverses the tree and replaces nodes with their majority class
// if pruning leads to a reduction in error.
func PruneTree(root *model.Node, validationSet *model.Dataset, targetAttr string) *model.Node {
	if root == nil {
		return nil // Return nil if the root is nil
	}

	// Recursively prune the left and right subtrees
	if root.Left != nil {
		root.Left = PruneTree(root.Left, validationSet, targetAttr)
	}
	if root.Right != nil {
		root.Right = PruneTree(root.Right, validationSet, targetAttr)
	}

	// Estimate the error of the current node
	errorWithoutPruning := EstimateError(root, validationSet, targetAttr)

	// Create a leaf node with the majority class of the validation set
	majorityClass := getClassification(validationSet, targetAttr) // Using the method here
	prunedNode := &model.Node{Class: majorityClass}

	// Estimate the error if we prune this node
	errorWithPruning := EstimateError(prunedNode, validationSet, targetAttr)

	// If pruning reduces the error, return the pruned node
	if errorWithPruning < errorWithoutPruning {
		return prunedNode
	}

	// Otherwise, return the current node
	return root
}

// EstimateError calculates the error rate of a node based on the validation dataset.
// It compares the predicted class of the node against the actual class in the dataset.
func EstimateError(node *model.Node, dataset *model.Dataset, targetAttr string) float64 {
	if node == nil {
		return 0.0
	}

	correctPredictions := 0
	totalPredictions := len(dataset.RowInstances)

	for _, instance := range dataset.RowInstances {
		predictedClass := node.Class
		actualClass := instance[targetAttr]

		if predictedClass == actualClass {
			correctPredictions++
		}
	}

	return 1.0 - float64(correctPredictions)/float64(totalPredictions)
}

// getMajorityClass returns the majority class from the dataset based on the target attribute.
// It counts the occurrences of each class and identifies the one with the highest count.
func getClassification(dataset *model.Dataset, targetAttr string) string {
	classCounts := make(map[string]int)
	for _, instance := range dataset.RowInstances {
		class := instance[targetAttr].(string)
		classCounts[class]++
	}

	var majorityClass string
	maxCount := 0
	for class, count := range classCounts {
		if count > maxCount {
			maxCount = count
			majorityClass = class
		}
	}

	return majorityClass
}
