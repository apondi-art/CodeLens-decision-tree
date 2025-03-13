package algorithm

import "CodeLens-decision-tree/internal/model"


// CalculateEntropy computes the entropy of a dataset based on a given target attribute.
// Entropy measures the level of impurity or uncertainty in the dataset.
//
// Formula:
// H(S) = - ∑ (p_i * log₂(p_i))
//
// where p_i is the probability of class i in the target attribute.
//
// Parameters:
// - dataset: A pointer to the dataset containing instances.
// - targetAttr: The name of the target attribute.
//
// Returns:
// - A float64 value representing the entropy of the dataset

func CalculateEntropy(dataset *model.Dataset, targetAttr string) float64
func CalculateInformationGain(dataset *model.Dataset, attr *model.Attribute, targetAttr string) float64
func CalculateGainRatio(dataset *model.Dataset, attr *model.Attribute, targetAttr string) float64
