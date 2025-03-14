package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"CodeLens-decision-tree/internal/algorithm"
	"CodeLens-decision-tree/internal/data"
	"CodeLens-decision-tree/internal/model"
	"CodeLens-decision-tree/internal/serialization"
)

// Flags holds the command-line arguments.
type Flags struct {
	Command      string
	InputFile    string
	TargetColumn string
	OutputPath   string
	TrainedFile  string
}

var CmdFlags Flags // Global variable for parsed flags

// init initializes command-line flags.
func init() {
	flag.StringVar(&CmdFlags.Command, "c", "", "Command: train or predict")
	flag.StringVar(&CmdFlags.InputFile, "i", "", "Path to input CSV file")
	flag.StringVar(&CmdFlags.TargetColumn, "t", "", "Name of target column (only for training)")
	flag.StringVar(&CmdFlags.OutputPath, "o", "", "Path to save output (model or predictions)")
	flag.StringVar(&CmdFlags.TrainedFile, "m", "", "Path to trained model file (only for prediction)")
}

// ParseFlags validates and returns the parsed flag values.
func ParseFlags() (string, string, string, string, string, error) {
	flag.Parse()

	if CmdFlags.Command != "train" && CmdFlags.Command != "predict" {
		return "", "", "", "", "", errors.New("invalid command: must be 'train' or 'predict'")
	}
	if CmdFlags.InputFile == "" {
		return "", "", "", "", "", errors.New("input file path is required")
	}
	if CmdFlags.OutputPath == "" {
		return "", "", "", "", "", errors.New("output path is required")
	}
	if CmdFlags.Command == "train" && CmdFlags.TargetColumn == "" {
		return "", "", "", "", "", errors.New("target column is required for training")
	}
	if CmdFlags.Command == "predict" && CmdFlags.TrainedFile == "" {
		return "", "", "", "", "", errors.New("trained model file is required for prediction")
	}

	return CmdFlags.Command, CmdFlags.InputFile, CmdFlags.TargetColumn, CmdFlags.TrainedFile, CmdFlags.OutputPath, nil
}

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

// ExecutePredictCommand loads a trained model and makes predictions.
func ExecutePredictCommand(inputPath, modelPath, outputPath string) error {
	// Load trained model.
	tree, err := serialization.DeserializeTree(modelPath)
	if err != nil {
		return fmt.Errorf("failed to load model: %v", err)
	}

	// Load input dataset.
	dataset, err := data.GenerateCSVData(inputPath, tree.TargetAttr)
	if err != nil {
		return fmt.Errorf("failed to load dataset: %v", err)
	}

	// Make predictions.
	predictions := make([]string, len(dataset.RowInstances))
	for i, instance := range dataset.RowInstances {
		predictions[i] = tree.Root.Predict(instance)
	}

	// Save predictions.
	err = savePredictions(predictions, outputPath)
	if err != nil {
		return fmt.Errorf("failed to save predictions: %v", err)
	}

	fmt.Println("Predictions completed successfully. Predictions saved to:", outputPath)
	return nil
}

// savePredictions writes predictions to a CSV file.
func savePredictions(predictions []string, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, prediction := range predictions {
		err := writer.Write([]string{prediction})
		if err != nil {
			return fmt.Errorf("failed to write prediction: %v", err)
		}
	}

	return nil
}

func main() {
	// Parse CLI flags.
	command, inputPath, targetColumn, modelPath, outputPath, err := ParseFlags()
	if err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	// Execute the appropriate command.
	switch command {
	case "train":
		err = ExecuteTrainCommand(inputPath, targetColumn, outputPath)
	case "predict":
		err = ExecutePredictCommand(inputPath, modelPath, outputPath)
	default:
		log.Fatalf("Invalid command: %s", command)
	}

	// Handle errors.
	if err != nil {
		log.Fatalf("Error executing command: %v", err)
	}

	fmt.Println("Command executed successfully.")
}
