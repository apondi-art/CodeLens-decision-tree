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
	for _, instance := range dataset.Records {
		if instance[attr.Name] == value {
			subset = append(subset, instance)
		}
	}

	return &model.Dataset{Records: subset, Attributes: dataset.Attributes}, nil
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
	for _, instance := range dataset.Records {
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
