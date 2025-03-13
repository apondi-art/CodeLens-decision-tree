package algorithm

import "CodeLens-decision-tree/internal/model"

func SplitDataset(dataset *model.Dataset, attr *model.Attribute, value interface{}) (*model.Dataset, error)
func FindBestNumericalSplit(dataset *model.Dataset, attr *model.Attribute, targetAttr string) (float64, float64)
func DistributeInstance(instance map[string]interface{}, children map[interface{}]*model.Node) map[*model.Node]float64
