package flags

import (
	"fmt"

	"CodeLens-decision-tree/internal/algorithm"
	"CodeLens-decision-tree/internal/data"
	"CodeLens-decision-tree/internal/model"
	"CodeLens-decision-tree/internal/serialization"
)

// executeTrainCommand trains a decision tree model and saves it.
func ExecuteTrainCommand(inputPath, targetColumn, outputPath string) error {
	// Load the dataset.
	dataset, err := data.GenerateCSVData(inputPath, targetColumn)
	if err != nil {
		return fmt.Errorf("failed to load dataset: %v", err)
	}

	// Validate dataset.
	if err := data.ValidateTargetColumn(dataset, targetColumn); err != nil {
		return fmt.Errorf("invalid target column: %v", err)
	}
	if err := data.ValidateDataCompleteness(dataset); err != nil {
		return fmt.Errorf("invalid dataset: %v", err)
	}

	// Extract attributes.
	attributes := extractAttributes(dataset)

	// Train decision tree.
	maxDepth := 10
	tree, err := algorithm.BuildTree(dataset, attributes, targetColumn, 0, maxDepth)
	if err != nil {
		return fmt.Errorf("failed to build decision tree: %v", err)
	}

	// Save model.
	if err := serialization.SerializeTree(tree, outputPath); err != nil {
		return fmt.Errorf("failed to save model: %v", err)
	}

	fmt.Println("Training completed successfully. Model saved to:", outputPath)
	return nil
}

// extractAttributes extracts attributes from the dataset.
func extractAttributes(dataset *model.Dataset) []*model.Attribute {
	var attributes []*model.Attribute
	for _, attr := range dataset.ColumnAttributes {
		attributes = append(attributes, attr)
	}
	return attributes
}
