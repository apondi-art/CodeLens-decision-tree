package model

import (
	"errors"
	"fmt"
	"math"
)

// Interface defines the methods that a dataset must implement.
type Interface interface {
	CalculateClassEntropy() float64
	GetUniqueValues(attributeName string) []interface{}
	IsPure() bool
	GetMajorityClass() string
	CountClassInstances() map[string]int
	SplitByCategoricalValue(attributeName string) (map[interface{}]*Dataset, error)
	SplitByNumericThreshold(attributeName string, threshold float64) (map[interface{}]*Dataset, error)
	FilterByAttributeValue(attributeName string, value interface{}) *Dataset
}

// DecisionTree represents a C4.5 decision tree model
type DecisionTree struct {
	Root            *TreeNode
	Attributes      []*Attribute
	TargetAttr      string
	AttributeTypes  map[string]AttributeType
	Metadata        map[string]interface{} // For storing training data stats, feature importance, etc.
	MaxDepth        int                    // Maximum depth used during training
	MinSamplesSplit int                    // Minimum samples required for a split
	PruningEnabled  bool                   // Whether pruning was applied
}

// SelectBestAttribute selects the attribute with the highest gain ratio.
func SelectBestAttribute(dataset Interface, attributes []*Attribute, targetAttr string) (*Attribute, float64) {
	var bestAttribute *Attribute
	maxGainRatio := -1.0

	for _, attr := range attributes {
		gainRatio := CalculateGainRatio(dataset, attr, targetAttr)
		fmt.Printf("Attribute: %s, Gain Ratio: %f\n", attr.Name, gainRatio) // Debugging Output

		if gainRatio > maxGainRatio {
			maxGainRatio = gainRatio
			bestAttribute = attr
		}
	}

	return bestAttribute, maxGainRatio
}

func convertCountsToFloat(counts map[string]int) map[string]float64 {
	result := make(map[string]float64)
	for class, count := range counts {
		result[class] = float64(count)
	}
	return result
}

// CalculateGainRatio calculates the gain ratio for a given attribute.
func CalculateGainRatio(dataset Interface, attribute *Attribute, targetAttr string) float64 {
	gain := CalculateInformationGain(dataset, attribute, targetAttr)
	splitInfo := CalculateSplitInfo(attribute, dataset)

	if splitInfo == 0 {
		return 0
	}

	gainRatio := gain / splitInfo
	return gainRatio
}

// CalculateInformationGain calculates the information gain for a given attribute.
func CalculateInformationGain(dataset Interface, attribute *Attribute, targetAttr string) float64 {
	entropy := dataset.CalculateClassEntropy()
	values := dataset.GetUniqueValues(attribute.Name)
	weightedEntropy := 0.0

	concreteDatasetForFilter, ok := dataset.(*Dataset)
	if !ok {
		return 0
	}

	for _, value := range values {
		concreteSubset := concreteDatasetForFilter.FilterByAttributeValue(attribute.Name, value)
		subsetProbability := float64(len(concreteSubset.RowInstances)) / float64(len(concreteDatasetForFilter.RowInstances))
		weightedEntropy += subsetProbability * concreteSubset.CalculateClassEntropy()
	}

	return entropy - weightedEntropy
}

// CalculateSplitInfo calculates the split information for a given attribute.
func CalculateSplitInfo(attribute *Attribute, dataset Interface) float64 {
	splitInfo := 0.0
	values := dataset.GetUniqueValues(attribute.Name)

	concreteDatasetForFilter, ok := dataset.(*Dataset)
	if !ok {
		return 0
	}

	for _, value := range values {
		concreteSubset := concreteDatasetForFilter.FilterByAttributeValue(attribute.Name, value)
		subsetProbability := float64(len(concreteSubset.RowInstances)) / float64(len(concreteDatasetForFilter.RowInstances))
		splitInfo -= subsetProbability * Log2(subsetProbability)
	}

	return splitInfo
}

// Log2 calculates the base-2 logarithm of a number.
func Log2(x float64) float64 {
	return math.Log(x) / math.Log(2)
}

// BuildTree constructs a decision tree using the C4.5 algorithm.
func (t *DecisionTree) Train(dataset Interface, attributes []*Attribute, targetAttr string, maxDepth int) error {
	if dataset == nil {
		return fmt.Errorf("dataset is empty or nil")
	}

	// Validate attributes and target attribute
	if len(attributes) == 0 || targetAttr == "" {
		return fmt.Errorf("attributes or target attribute are not set")
	}

	// Build the decision tree using the algorithm package
	root, err := BuildTree(dataset, attributes, targetAttr, 0, maxDepth)
	if err != nil {
		return fmt.Errorf("failed to build decision tree: %v", err)
	}

	t.Root = root // Set the root node of the DecisionTree

	return nil
}

// Predict classifies a single instance
func (t *DecisionTree) Predict(instance map[string]interface{}) string {
	if t.Root == nil {
		return ""
	}
	return t.Root.Predict(instance)
}

// BatchPredict classifies multiple instances
func (t *DecisionTree) BatchPredict(dataset *Dataset) []string {
	if t.Root == nil || dataset == nil {
		return nil
	}

	predictions := make([]string, len(dataset.RowInstances))
	for i, instance := range dataset.RowInstances {
		predictions[i] = t.Predict(instance)
	}
	return predictions
}

// PredictWithConfidence returns predicted class and confidence score
func (t *DecisionTree) PredictWithConfidence(instance map[string]interface{}) (string, float64) {
	if t.Root == nil {
		return "", 0
	}

	// Traverse the tree to get the predicted class and class distribution
	// Traverse the tree to get the predicted class and class distribution
	predictedClass, classDistribution := t.Root.PredictWithDistribution(instance)

	// Calculate confidence as the proportion of the predicted class in the distribution
	total := 0.0
	for _, count := range classDistribution {
		total += count
	}
	confidence := classDistribution[predictedClass] / total

	return predictedClass, confidence
}

// GetFeatureImportance returns importance score for each attribute
func (t *DecisionTree) GetFeatureImportance() map[string]float64 {
	if t.Metadata == nil {
		return nil
	}

	// Assume Metadata contains feature importance scores
	if importance, ok := t.Metadata["featureImportance"].(map[string]float64); ok {
		return importance
	}
	return nil
}

// TreeNode represents a single node in the decision tree
type TreeNode struct {
	SplitAttr         string
	SplitValue        interface{}
	IsLeaf            bool
	PredictedClass    string
	ClassDistribution map[string]float64
	Children          map[string]*TreeNode
}

// PredictWithDistribution traverses the tree and returns the predicted class and class distribution

// Predict traverses the tree to classify a single instance
func (n *TreeNode) Predict(instance map[string]interface{}) string {
	if n.IsLeaf {
		return n.PredictedClass
	}

	// Get the value of the split attribute for the instance
	value, exists := instance[n.SplitAttr]
	if !exists {
		return "" // Return empty string if attribute is missing
	}

	// Handle numerical splits
	if threshold, ok := n.SplitValue.(float64); ok {
		if val, ok := value.(float64); ok && val <= threshold {
			return n.Children["left"].Predict(instance)
		} else if ok {
			return n.Children["right"].Predict(instance)
		}
	}

	// Handle categorical splits
	if child, ok := n.Children[fmt.Sprintf("%v", value)]; ok {
		return child.Predict(instance)
	}

	return n.PredictedClass // Default to the predicted class if no match
}

// PredictWithDistribution traverses the tree and returns the predicted class and class distribution
func (n *TreeNode) PredictWithDistribution(instance map[string]interface{}) (string, map[string]float64) {
	if n.IsLeaf {
		return n.PredictedClass, n.ClassDistribution
	}

	// Get the value of the split attribute for the instance
	value, exists := instance[n.SplitAttr]
	if !exists {
		return "", n.ClassDistribution
	}

	// Handle numerical splits
	if threshold, ok := n.SplitValue.(float64); ok {
		if val, ok := value.(float64); ok && val <= threshold {
			return n.Children["left"].PredictWithDistribution(instance)
		} else if ok {
			return n.Children["right"].PredictWithDistribution(instance)
		}
	}

	// Handle categorical splits
	if child, ok := n.Children[fmt.Sprintf("%v", value)]; ok {
		return child.PredictWithDistribution(instance)
	}

	return n.PredictedClass, n.ClassDistribution
}

// BuildTree constructs a decision tree using the C4.5 algorithm.
func BuildTree(dataset Interface, attributes []*Attribute, targetAttr string, depth int, maxDepth int) (*TreeNode, error) {
	// Check if the maximum depth has been reached.
	if dataset == nil {
		return nil, fmt.Errorf("dataset is empty or nil")
	}

	if len(dataset.GetUniqueValues(targetAttr)) == 0 {
		return nil, errors.New("dataset is empty")
	}

	// Base case: If the dataset is pure, no attributes left, or max depth reached, return a leaf node.
	if dataset.IsPure() || len(attributes) == 0 || depth >= maxDepth {
		majorityClass := dataset.GetMajorityClass()
		classCounts := dataset.CountClassInstances()
		classDistribution := make(map[string]float64)
		for class, count := range classCounts {
			classDistribution[class] = float64(count)
		}
		return &TreeNode{
			IsLeaf:            true,
			PredictedClass:    majorityClass,
			ClassDistribution: classDistribution,
		}, nil
	}

	// Select the best attribute to split on.
	datasetPtr, ok := dataset.(*Dataset)
	if !ok {
		return nil, errors.New("dataset is not of type *Dataset")
	}
	bestAttribute, _ := SelectBestAttribute(datasetPtr, attributes, targetAttr)
	if bestAttribute == nil {
		return nil, errors.New("no best attribute found")
	}

	// Create a decision node for the best attribute.
	node := &TreeNode{
		SplitAttr:         bestAttribute.Name,
		IsLeaf:            false,
		ClassDistribution: convertCountsToFloat(dataset.CountClassInstances()),
	}

	// Remove the best attribute from the list of remaining attributes.
	var remainingAttributes []*Attribute
	for _, attr := range attributes {
		if attr.Name != bestAttribute.Name {
			remainingAttributes = append(remainingAttributes, attr)
		}
	}

	// Split the dataset based on the best attribute.
	var subsets map[interface{}]*Dataset
	var err error

	if bestAttribute.Type == Categorical {
		subsets, err = dataset.SplitByCategoricalValue(bestAttribute.Name)
	} else if bestAttribute.Type == Numeric {
		// For numerical attributes, find the best split threshold.
		datasetPtr, ok := dataset.(*Dataset)
		if !ok {
			return nil, errors.New("dataset is not of type *Dataset")
		}
		split, _ := bestAttribute.FindBestSplit(datasetPtr)
		node.SplitValue = split.Value
		subsets, err = dataset.SplitByNumericThreshold(bestAttribute.Name, split.Value.(float64))
	}

	if err != nil {
		return nil, err
	}

	node.Children = make(map[string]*TreeNode)
	for value, subset := range subsets {
		if len(subset.RowInstances) == 0 {
			node.Children[fmt.Sprintf("%v", value)] = &TreeNode{
				IsLeaf:            true,
				PredictedClass:    dataset.GetMajorityClass(),
				ClassDistribution: convertCountsToFloat(dataset.CountClassInstances()),
			}
			continue
		}
		childNode, err := BuildTree(subset, remainingAttributes, targetAttr, depth+1, maxDepth)
		if err != nil {
			return nil, err
		}
		node.Children[fmt.Sprintf("%v", value)] = childNode
	}

	return node, nil
}
