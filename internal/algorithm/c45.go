package algorithm

import (
	"fmt"

	"CodeLens-decision-tree/internal/model"
)

// SelectBestAttribute selects the attribute with the highest gain ratio.
func SelectBestAttribute(dataset model.DatasetInterface, attributes []*model.Attribute, targetAttr string) (*model.Attribute, float64) {
	var bestAttribute *model.Attribute
	maxGainRatio := -1.0

	datasetPtr, ok := dataset.(*model.Dataset)
	if !ok {
		fmt.Println("Failed type assertion")
		return nil, maxGainRatio
	}
	for _, attr := range attributes {
		if CalculateGainRatio(datasetPtr, attr, targetAttr) > maxGainRatio {
			maxGainRatio = CalculateGainRatio(datasetPtr, attr, targetAttr)
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
