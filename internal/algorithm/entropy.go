package algorithm

import (
	"math"

	"CodeLens-decision-tree/internal/model"
)

// CalculateEntropy computes the entropy of a dataset based on a given target attribute.
// Entropy measures the level of impurity or uncertainty in the dataset.
//
// Formula:
// H(S) = - ∑ (p_i * log₂(p_i))
//
// where p_i is the probability of class i in the target attribute.
//
// Parameters:
// - dataset: A pointer to the dataset containing instances.
// - targetAttr: The name of the target attribute.
//
// Returns:
// - A float64 value representing the entropy of the dataset

func CalculateEntropy(dataset *model.Dataset, targetAttr string) float64 {
	if len(dataset.RowInstances) == 0 {
		return 0.0
	}

	// Use cached values if target column matches
	var classCounts map[string]int
	if dataset.TargetColumn == targetAttr && dataset.TargetOccurrence != nil {
		classCounts = dataset.TargetOccurrence
	} else {
		classCounts = make(map[string]int)
		for _, instance := range dataset.RowInstances {
			if classValue, ok := instance[targetAttr].(string); ok && classValue != "" {
				classCounts[classValue]++
			}
		}

		// Cache results for reuse
		if dataset.TargetColumn == targetAttr {
			dataset.TargetOccurrence = classCounts
		}
	}

	totalInstances := len(dataset.RowInstances)

	// Quick return for pure datasets
	if len(classCounts) <= 1 {
		return 0.0
	}

	var entropy float64
	for _, count := range classCounts {
		if count > 0 {
			probability := float64(count) / float64(totalInstances)
			entropy -= probability * getLog2(probability)
		}
	}

	return entropy
}

// getLog2 retrieves or calculates log2 value with caching
func getLog2(value float64) float64 {
	// Check cache first with read lock
	log2CacheMutex.RLock()
	if cached, exists := log2Cache[value]; exists {
		log2CacheMutex.RUnlock()
		return cached
	}
	log2CacheMutex.RUnlock()

	// Calculate and cache if common value

	result := math.Log2(value)

	// Only cache common values to prevent memory bloat
	if value > 0.01 && value < 1.0 {
		log2CacheMutex.Lock()
		log2Cache[value] = result
		log2CacheMutex.Unlock()
	}

	return result
}

// CalculateGainRatio computes the gain ratio of splitting a dataset on a given attribute.
// Gain ratio is an improvement over information gain that normalizes for the number of possible splits.
//
// Formula:
// GainRatio(S, A) = IG(S, A) / SplitInfo(S, A)
//
// where:
// - IG(S, A) is the information gain of splitting on attribute A,
// - SplitInfo(S, A) = - ∑ (|S_v| / |S|) * log₂(|S_v| / |S|)
//   is the entropy of the distribution of instances across subsets.
//
// Parameters:
// - dataset: A pointer to the dataset containing instances.
// - attr: The attribute used for splitting.
// - targetAttr: The name of the target attribute.
//
// Returns:
// - A float64 value representing the gain ratio.

func CalculateInformationGain(dataset *model.Dataset, attr *model.Attribute, targetAttr string) float64 {
	originalEntropy := CalculateEntropy(dataset, targetAttr)
	if originalEntropy == 0 {
		return 0.0
	}

	subsets := make(map[interface{}]*model.Dataset)
	for _, instance := range dataset.RowInstances {
		attrValue := instance[attr.Name]
		if attrValue == nil { // Skip instances with missing attribute values
			continue
		}

		// Create a Split struct
		split := &model.Split{
			Attribute: attr,
			Type:      attr.Type.String(), // Convert AttributeType to string
			Value:     attrValue,
		}

		// Handle categorical splits by adding a map
		if attr.Type == model.Categorical {
			split.CategoricalMap = map[string]bool{attrValue.(string): true}
		}

		subset, err := SplitDataset(dataset, split)
		if err != nil {
			continue // Handle error appropriately
		}
		subsets[attrValue] = subset
	}

	// Compute weighted entropy
	var weightedEntropy float64
	totalInstances := float64(len(dataset.RowInstances))
	for _, subDataset := range subsets {
		subsetProbability := float64(len(subDataset.RowInstances)) / totalInstances
		weightedEntropy += subsetProbability * CalculateEntropy(subDataset, targetAttr)
	}

	// Round result to avoid floating-point precision errors
	return math.Round((originalEntropy-weightedEntropy)*1e6) / 1e6
}

// CalculateGainRatio computes the gain ratio of splitting a dataset on a given attribute.
// Gain ratio is an improvement over information gain that normalizes for the number of possible splits.
//
// Formula:
// GainRatio(S, A) = IG(S, A) / SplitInfo(S, A)
//
// where:
// - IG(S, A) is the information gain of splitting on attribute A,
// - SplitInfo(S, A) = - ∑ (|S_v| / |S|) * log₂(|S_v| / |S|)
//   is the entropy of the distribution of instances across subsets.
//
// Parameters:
// - dataset: A pointer to the dataset containing instances.
// - attr: The attribute used for splitting.
// - targetAttr: The name of the target attribute.
//
// Returns:
// - A float64 value representing the gain ratio.

func CalculateGainRatio(dataset *model.Dataset, attr *model.Attribute, targetAttr string) float64 {
	informationGain := CalculateInformationGain(dataset, attr, targetAttr)
	if informationGain == 0 {
		return 0.0
	}

	splitInfo := 0.0
	totalInstances := float64(len(dataset.RowInstances))

	// Pre-allocate with reasonable capacity
	subsets := make(map[interface{}]int, 10)

	for _, instance := range dataset.RowInstances {
		attrValue := instance[attr.Name]
		if attrValue == nil { // Skip instances with missing attribute values
			continue
		}
		subsets[attrValue]++
	}

	for _, count := range subsets {
		probability := float64(count) / totalInstances
		if probability > 0 {
			splitInfo -= probability * getLog2(probability)
		}
	}

	if splitInfo == 0 {
		return 0.0
	}

	return informationGain / splitInfo
}
