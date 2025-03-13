package algorithm

import "CodeLens-decision-tree/internal/model"

func CalculateEntropy(dataset *model.Dataset, targetAttr string) float64
func CalculateInformationGain(dataset *model.Dataset, attr *model.Attribute, targetAttr string) float64
func CalculateGainRatio(dataset *model.Dataset, attr *model.Attribute, targetAttr string) float64
