package algorithm

import( 
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
	if len(dataset.Instances) == 0 {
		return 0.0
	}

	// Count occurrences of each class in the target attribute.
	classCounts := make(map[string]int)
	for _, instance := range dataset.Instances {
		classValue, ok := instance[targetAttr].(string)
		if !ok {
			continue // Ignore instances with missing or non-string target values.
		}
		classCounts[classValue]++
	}

	// Calculate entropy
	var entropy float64
	totalInstances := float64(len(dataset.Instances))
	for _, count := range classCounts {
		probability := float64(count) / totalInstances
		entropy -= probability * math.Log2(probability)
	}

	return entropy
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
	// Compute the entropy before splitting
	originalEntropy := CalculateEntropy(dataset, targetAttr)

	// Group instances by attribute value
	subsets := make(map[interface{}][]map[string]interface{})
	for _, instance := range dataset.Instances {
		attrValue := instance[attr.Name]
		subsets[attrValue] = append(subsets[attrValue], instance)
	}

	// Compute weighted entropy after splitting
	var weightedEntropy float64
	totalInstances := float64(len(dataset.Instances))
	for _, subset := range subsets {
		subDataset := &model.Dataset{Instances: subset}
		probability := float64(len(subDataset.Instances)) / totalInstances
		weightedEntropy += probability * CalculateEntropy(subDataset, targetAttr)
	}

	// Information Gain = Entropy before split - Weighted entropy after split
	return originalEntropy - weightedEntropy
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
	// Compute the information gain
	informationGain := CalculateInformationGain(dataset, attr, targetAttr)
	if informationGain == 0 {
		return 0.0 // Avoid division by zero
	}

	// Compute the SplitInfo (entropy of the attribute distribution)
	var splitInfo float64
	totalInstances := float64(len(dataset.Instances))
	subsets := make(map[interface{}]int)

	// Count instances per unique attribute value
	for _, instance := range dataset.Instances {
		attrValue := instance[attr.Name]
		subsets[attrValue]++
	}

	// Compute SplitInfo
	for _, count := range subsets {
		probability := float64(count) / totalInstances
		splitInfo -= probability * math.Log2(probability)
	}

	// Gain Ratio = Information Gain / SplitInfo
	if splitInfo == 0 {
		return 0.0 // Prevent division by zero
	}
	return informationGain / splitInfo
}
