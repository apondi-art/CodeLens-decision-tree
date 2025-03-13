package algorithm

import "CodeLens-decision-tree/internal/model"

func BuildTree(dataset *model.Dataset, attributes []*model.Attribute, targetAttr string, depth int, maxDepth int) (*model.Node, error)
func SelectBestAttribute(dataset *model.Dataset, attributes []*model.Attribute, targetAttr string) (*model.Attribute, float64)
func MajorityClass(dataset *model.Dataset, targetAttr string) string
