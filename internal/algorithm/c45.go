package algorithm

import (
	"errors"

	"CodeLens-decision-tree/internal/model"
)

// BuildTree constructs a decision tree using the C4.5 algorithm.
func BuildTree(dataset *model.Dataset, attributes []*model.Attribute, targetAttr string, depth int, maxDepth int) (*model.Node, error) {
	// Check for empty dataset.
	if len(dataset.RowInstances) == 0 {
		return nil, errors.New("dataset is empty")
	}
	
	// Base case: If the dataset is pure or no attributes are left, return a leaf node.
	if dataset.IsPure() || len(attributes) == 0 || depth >= maxDepth {
		majorityClass := dataset.GetMajorityClass()
		return &model.Node{
			IsLeaf:            true,
			PredictedClass:    majorityClass,
			Depth:             depth,
			SampleCount:       len(dataset.RowInstances),
			ClassDistribution: dataset.CountClassInstances(),
		}, nil
	}

	// Select the best attribute to split on.
	bestAttribute, gainRatio := SelectBestAttribute(*dataset, attributes, targetAttr)
	if bestAttribute == nil {
		return nil, errors.New("no best attribute found")
	}

	// Create a decision node for the best attribute.
	node := &model.Node{
		Attribute:         bestAttribute,
		IsLeaf:            false,
		Depth:             depth,
		SampleCount:       len(dataset.RowInstances),
		ClassDistribution: dataset.CountClassInstances(),
		GainRatio:         gainRatio,
		Children:          make(map[interface{}]*model.Node),
	}

	// Remove the best attribute from the list of remaining attributes.
	var remainingAttributes []*model.Attribute
	for _, attr := range attributes {
		if attr.Name != bestAttribute.Name {
			remainingAttributes = append(remainingAttributes, attr)
		}
	}

	// Split the dataset based on the best attribute.
	var subsets map[interface{}]*model.Dataset
	var err error

	if bestAttribute.Type == model.Categorical {
		subsets, err = dataset.SplitByCategoricalValue(bestAttribute.Name)
	} else if bestAttribute.Type == model.Numeric {
		// For numerical attributes, find the best split threshold.
		split, _ := bestAttribute.FindBestSplit(dataset)
		subsets, err = dataset.SplitByNumericThreshold(bestAttribute.Name, split.Value.(float64))
	}

	if err != nil {
		return nil, err
	}

	// Recursively build the tree for each subset.
	for value, subset := range subsets {
		childNode, err := BuildTree(subset, remainingAttributes, targetAttr, depth+1, maxDepth)
		if err != nil {
			return nil, err
		}
		node.Children[value] = childNode
	}

	return node, nil
}

// SelectBestAttribute selects the attribute with the highest gain ratio.
func SelectBestAttribute(dataset model.Dataset, attributes []*model.Attribute, targetAttr string) (*model.Attribute, float64) {
	var bestAttribute *model.Attribute
	maxGainRatio := -1.0

	for _, attr := range attributes {
		gainRatio := CalculateGainRatio(&dataset, attr, targetAttr)
		if gainRatio > maxGainRatio {
			maxGainRatio = gainRatio
			bestAttribute = attr
		}
	}

	return bestAttribute, maxGainRatio
}

// MajorityClass returns the most frequent class in the dataset for the target attribute.
func MajorityClass(dataset *model.Dataset, targetAttr string) string {
	classCounts := dataset.CountClassInstances()
	majorityClass := ""
	maxCount := 0

	for class, count := range classCounts {
		if count > maxCount {
			majorityClass = class
			maxCount = count
		}
	}

	return majorityClass
}
