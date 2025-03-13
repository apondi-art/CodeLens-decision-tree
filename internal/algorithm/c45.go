package algorithm

import "CodeLens-decision-tree/internal/model"

func BuildTree(dataset *model.Dataset, attributes []*model.Attribute, targetAttr string, depth int, maxDepth int) (*model.Node, error) {
	return &model.Node{}, nil
}
func SelectBestAttribute(dataset *model.Dataset, attributes []*model.Attribute, targetAttr string) (*model.Attribute, float64) {
	return &model.Attribute{}, 0
}
func MajorityClass(dataset *model.Dataset, targetAttr string) string {
	return ""
}
