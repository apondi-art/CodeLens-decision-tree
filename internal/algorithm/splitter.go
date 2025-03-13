package algorithm

import (
	"errors"
	"sort"

	"CodeLens-decision-tree/internal/model"
)

// SplitDataset divides the dataset based on a specified attribute value.
// It returns a subset where the attribute matches the provided value.
//
// Parameters:
// - dataset: A pointer to the dataset to be split.
// - attr: A pointer to the attribute used for splitting.
// - value: The attribute value to filter records by.
//
// Returns:
// - A new dataset containing only records matching the given attribute value.
// - An error if dataset or attribute is nil.
func SplitDataset(dataset *model.Dataset, attr *model.Attribute, value interface{}) (*model.Dataset, error) {
	if dataset == nil || attr == nil {
		return nil, errors.New("dataset or attribute is nil")
	}

	var subset []map[string]interface{}
	for _, instance := range dataset.RowInstances {
		if instance[attr.Name] == value {
			subset = append(subset, instance)
		}
	}

	return &model.Dataset{RowInstances: subset, ColumnAttributes: dataset.ColumnAttributes}, nil
}

// FindBestNumericalSplit identifies the optimal threshold for splitting a numerical attribute.
// It evaluates potential thresholds using information gain and selects the best one.
//
// Parameters:
// - dataset: A pointer to the dataset containing numerical data.
// - attr: A pointer to the numerical attribute to evaluate.
// - targetAttr: The target attribute for calculating information gain.
//
// Returns:
// - The optimal threshold value for splitting.
// - The corresponding information gain.
func FindBestNumericalSplit(dataset *model.Dataset, attr *model.Attribute, targetAttr string) (float64, float64) {
	if dataset == nil || attr == nil {
		return 0, 0
	}

	var bestThreshold, bestGain float64
	var values []float64

	// Extract numerical values from dataset
	for _, instance := range dataset.RowInstances {
		if val, ok := instance[attr.Name].(float64); ok {
			values = append(values, val)
		}
	}

	// Sort values to determine possible split points
	sort.Float64s(values)

	// Iterate over midpoints of consecutive values to find the best split
	for i := 1; i < len(values); i++ {
		threshold := (values[i-1] + values[i]) / 2
		gain := computeInformationGain(dataset, attr.Name, threshold, targetAttr)
		if gain > bestGain {
			bestGain = gain
			bestThreshold = threshold
		}
	}

	return bestThreshold, bestGain
}

// DistributeInstance assigns an instance to child nodes based on attribute values.
// This function is useful for decision trees where a probabilistic split might occur.
//
// Parameters:
// - instance: A map representing the data instance.
// - children: A map of possible child nodes indexed by attribute values.
//
// Returns:
// - A mapping of child nodes to their respective probabilities.
func DistributeInstance(instance map[string]interface{}, children map[interface{}]*model.Node) map[*model.Node]float64 {
	distribution := make(map[*model.Node]float64)

	for value, child := range children {
		if instanceValue, exists := instance[child.Attribute.Name]; exists && instanceValue == value {
			distribution[child] = 1.0
		}
	}

	return distribution
}

// computeInformationGain estimates the information gain for a numerical split.
// This function serves as a placeholder and should be replaced with an actual entropy-based calculation.
//
// Parameters:
// - dataset: The dataset containing instances.
// - attr: The attribute to evaluate.
// - threshold: The numerical split point.
// - targetAttr: The target attribute for classification.
//
// Returns:
// - A float64 representing the computed information gain.
func computeInformationGain(dataset *model.Dataset, attr string, threshold float64, targetAttr string) float64 {
	// Calculate original entropy
	originalEntropy := CalculateEntropy(dataset, targetAttr)

	// Create subsets based on the threshold
	subset1 := &model.Dataset{RowInstances: []map[string]interface{}{}}
	subset2 := &model.Dataset{RowInstances: []map[string]interface{}{}}

	for _, instance := range dataset.RowInstances {
		if instance[attr].(float64) <= threshold {
			subset1.RowInstances = append(subset1.RowInstances, instance)
		} else {
			subset2.RowInstances = append(subset2.RowInstances, instance)
		}
	}

	// Calculate weighted entropy
	totalInstances := float64(len(dataset.RowInstances))
	weightedEntropy := (float64(len(subset1.RowInstances))/totalInstances)*CalculateEntropy(subset1, targetAttr) +
		(float64(len(subset2.RowInstances))/totalInstances)*CalculateEntropy(subset2, targetAttr)

	// Return information gain
	return originalEntropy - weightedEntropy
}
