package algorithm

import "CodeLens-decision-tree/internal/model"

func PruneTree(root *model.Node, validationSet *model.Dataset, targetAttr string) *model.Node
func EstimateError(node *model.Node, dataset *model.Dataset, targetAttr string) float64
