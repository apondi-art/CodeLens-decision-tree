package flags

import (
	"encoding/csv"
	"fmt"
	"os"

	"CodeLens-decision-tree/internal/data"
	"CodeLens-decision-tree/internal/serialization"
)

// ExecutePredictCommand loads a trained model and makes predictions.
func ExecutePredictCommand(inputPath, modelPath, outputPath string) error {
	// Load trained model
	tree, err := serialization.DeserializeTree(modelPath)
	if err != nil {
		return fmt.Errorf("failed to load model: %v", err)
	}

	// Load input dataset
	dataset, err := data.GenerateCSVData(inputPath, tree.TargetAttr)
	if err != nil {
		return fmt.Errorf("failed to load dataset: %v", err)
	}

	// Make predictions
	predictions := tree.BatchPredict(dataset)

	// Optionally, calculate confidence scores
	confidenceScores := make([]float64, len(dataset.RowInstances))
	for i, instance := range dataset.RowInstances {
		_, confidence := tree.PredictWithConfidence(instance)
		confidenceScores[i] = confidence
	}

	// Save predictions
	err = savePredictionsWithConfidence(predictions, confidenceScores, outputPath)
	if err != nil {
		return fmt.Errorf("failed to save predictions: %v", err)
	}

	fmt.Println("Predictions completed successfully. Predictions saved to:", outputPath)

	// Print feature importance
	featureImportance := tree.GetFeatureImportance()
	fmt.Println("\nFeature Importance:")
	for attr, importance := range featureImportance {
		fmt.Printf("%s: %.4f\n", attr, importance)
	}

	return nil
}

// savePredictionsWithConfidence writes predictions and confidence scores to a CSV file.
func savePredictionsWithConfidence(predictions []string, confidenceScores []float64, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	headers := []string{"ID", "Predicted", "Confidence"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write header: %v", err)
	}

	// Write predictions with confidence scores
	for i, prediction := range predictions {
		row := []string{
			fmt.Sprintf("%d", i+1),
			prediction,
			fmt.Sprintf("%.4f", confidenceScores[i]),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write prediction: %v", err)
		}
	}

	return nil
}
