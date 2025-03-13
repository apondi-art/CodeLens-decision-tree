package algorithm

import "CodeLens-decision-tree/internal/model"

func BuildTree(dataset *model.Dataset, attributes []*model.Attribute, targetAttr string, depth int, maxDepth int) (*model.Node, error) {
	return &model.Node{}, nil
}

func SelectBestAttribute(dataset *model.Dataset, attributes []*model.Attribute, targetAttr string) (*model.Attribute, float64) {
	return &model.Attribute{}, 0
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
