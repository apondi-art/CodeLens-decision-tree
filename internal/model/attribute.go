package model

import (
	"fmt"
	"math"
	"sort"
)

type AttributeType int

const (
	Unknown AttributeType = iota
	Numeric
	Categorical
	Boolean
	Timestamp
)

// string method for AttributeType
func (a AttributeType) String() string {
	switch a {

	case Numeric:
		return "numeric"
	case Categorical:
		return "categorical"
	case Boolean:
		return "boolean"
	case Timestamp:
		return "timestamp"
	default:
		return "unknown"
	}
}

// Split represents a decision split point for an attribute
// Split represents a decision split point for an attribute
type Split struct {
	Attribute      *Attribute      // The attribute being split on
	Type           string          // "categorical" or "numerical"
	Value          interface{}     // Threshold value for numerical splits
	CategoricalMap map[string]bool // Map of categorical values for categorical splits
	GainRatio      float64         // The gain ratio of this split
}

type Attribute struct {
	Name           string
	Type           AttributeType
	PossibleValues []interface{} // For categorical attributes
	Min            float64       // For numerical attributes
	Max            float64       // For numerical attributes
	MissingCount   int           // Count of missing values
	// Additional fields as needed
}

// / CalculateGainRatio calculates the gain ratio for an attribute
func (a *Attribute) CalculateGainRatio(dataset *Dataset) float64 {
	if dataset.TotalRows == 0 {
		return 0
	}

	// Calculate the entropy of the class distribution before the split
	classEntropy := dataset.CalculateClassEntropy()
	if classEntropy == 0 {
		return 0 // If the dataset is pure, there's no information gain
	}

	var subsets map[interface{}]*Dataset
	var err error

	// Get the subsets based on the attribute type
	if a.Type == Numeric {
		// For numeric attributes, find the best split point
		threshold, _ := a.findBestNumericSplit(dataset)
		if threshold != nil {
			subsets, err = dataset.SplitByNumericThreshold(a.Name, *threshold)
		}
	} else {
		// For categorical attributes, split by each value
		subsets, err = dataset.SplitByCategoricalValue(a.Name)
	}

	if err != nil || len(subsets) <= 1 {
		return 0
	}

	// Prepare slices for gonum operations
	weights := make([]float64, 0, len(subsets))
	entropies := make([]float64, 0, len(subsets))

	totalRows := float64(dataset.TotalRows)

	// Calculate weights and entropies for each subset
	for _, subset := range subsets {
		if subset.TotalRows == 0 {
			continue
		}
		weight := float64(subset.TotalRows) / totalRows
		weights = append(weights, weight)
		entropies = append(entropies, subset.CalculateClassEntropy())
	}

	// Calculate weighted entropy
	weightedEntropy := 0.0
	for i := range weights {
		weightedEntropy += weights[i] * entropies[i]
	}

	// Calculate split information for gain ratio
	splitInfo := 0.0
	for _, weight := range weights {
		if weight > 0 {
			splitInfo -= weight * math.Log2(weight)
		}
	}

	// Calculate information gain
	infoGain := classEntropy - weightedEntropy

	// Avoid division by zero
	if splitInfo == 0 {
		return 0
	}

	// Calculate gain ratio
	gainRatio := infoGain / splitInfo
	return gainRatio
}

// FindBestSplit finds the best split for an attribute
func (a *Attribute) FindBestSplit(dataset *Dataset) (Split, float64) {
	split := Split{
		Attribute: a,
		GainRatio: 0,
	}

	// Early return for empty datasets
	if dataset == nil || len(dataset.RowInstances) == 0 {
		return split, 0
	}

	gainRatio := a.CalculateGainRatio(dataset)

	if a.Type == Numeric {
		// For numeric attributes, find the best split point
		threshold, err := a.findBestNumericSplit(dataset)
		if err == nil && threshold != nil {
			split.Type = "numerical"
			split.Value = *threshold
			split.GainRatio = gainRatio
		}
	} else {
		// For categorical attributes
		split.Type = "categorical"
		split.CategoricalMap = make(map[string]bool)

		// Get all unique values for the attribute
		uniqueValues := dataset.GetUniqueValues(a.Name)

		// Use pre-allocation for better performance
		split.CategoricalMap = make(map[string]bool, len(uniqueValues))
		for _, val := range uniqueValues {
			if strVal, ok := val.(string); ok {
				split.CategoricalMap[strVal] = true
			}
		}

		split.GainRatio = gainRatio
	}

	return split, gainRatio
}

// findBestNumericSplit finds the best threshold for a numeric attribute
func (a *Attribute) findBestNumericSplit(dataset *Dataset) (*float64, error) {
	if a.Type != Numeric {
		return nil, fmt.Errorf("attribute '%s' is not numeric", a.Name)
	}

	values, err := dataset.GetNumericValues(a.Name)
	if err != nil || len(values) <= 1 {
		return nil, fmt.Errorf("insufficient numeric values for attribute '%s'", a.Name)
	}

	// Sort the values
	sort.Float64s(values)

	// Pre-allocate with reasonable capacity
	estimatedSplitPoints := len(values) / 2
	splitPoints := make([]float64, 0, estimatedSplitPoints)

	// Find potential split points (midpoints between adjacent distinct values)
	for i := 0; i < len(values)-1; i++ {
		if values[i] != values[i+1] {
			splitPoints = append(splitPoints, (values[i]+values[i+1])/2)
		}
	}

	if len(splitPoints) == 0 {
		return nil, fmt.Errorf("no valid split points for attribute '%s'", a.Name)
	}

	// For large datasets with many split points, sample them
	if len(splitPoints) > 100 {
		// Take every nth point
		n := len(splitPoints) / 100
		sampledPoints := make([]float64, 0, 100)
		for i := 0; i < len(splitPoints); i += n {
			sampledPoints = append(sampledPoints, splitPoints[i])
		}
		splitPoints = sampledPoints
	}

	// Evaluate each split point
	bestGainRatio := -1.0
	var bestThreshold *float64

	// Cache the class entropy calculation
	classEntropy := dataset.CalculateClassEntropy()

	for _, threshold := range splitPoints {
		// Split the dataset
		subsets, err := dataset.SplitByNumericThreshold(a.Name, threshold)
		if err != nil || len(subsets) <= 1 {
			continue
		}
		// Calculate weighted entropy after split
		weightedEntropy := 0.0
		splitInfo := 0.0

		for _, subset := range subsets {
			if subset.TotalRows == 0 {
				continue
			}

			weight := float64(subset.TotalRows) / float64(dataset.TotalRows)
			weightedEntropy += weight * subset.CalculateClassEntropy()

			// Calculate split information for gain ratio
			splitInfo -= weight * math.Log2(weight)
		}

		// Calculate information gain
		infoGain := classEntropy - weightedEntropy

		// Avoid division by zero
		if splitInfo == 0 {
			continue
		}

		// Calculate gain ratio
		gainRatio := infoGain / splitInfo

		if gainRatio > bestGainRatio {
			bestGainRatio = gainRatio
			thresholdCopy := threshold
			bestThreshold = &thresholdCopy
		}
	}

	return bestThreshold, nil
}
