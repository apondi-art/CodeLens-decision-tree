package algorithm

import "CodeLens-decision-tree/internal/model"

func BuildTree(dataset *model.Dataset, attributes []*model.Attribute, targetAttr string, depth int, maxDepth int) (*model.Node, error) {
	return &model.Node{}, nil
}

// SelectBestAttribute selects the attribute with the highest gain ratio.
func SelectBestAttribute(dataset model.Dataset, attributes []model.Attribute, targetAttr string) (*model.Attribute, float64) {
	var bestAttribute *model.Attribute
	maxGainRatio := -1.0

	for _, attr := range attributes {
		gainRatio := CalculateGainRatio(&dataset, &attr, targetAttr)
		if gainRatio > maxGainRatio {
			maxGainRatio = gainRatio
			bestAttribute = &attr
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
