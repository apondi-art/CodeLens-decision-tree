package main

import (
	"flag"
)

// Flags holds the command-line arguments.
type Flags struct {
	Command      string // Specifies the action: "train" or "predict"
	InputFile    string // Path to the input CSV file (training or prediction data)
	TargetColumn string // Column containing target labels (only for training)
	OutputPath   string // Path to save the output file (model or predictions)
	TrainedFile  string // Path to the trained decision tree model file (used for prediction)
}

var cmdFlags Flags // Global variable to store parsed flags

// init initializes the command-line flags
func init() {
	flag.StringVar(&cmdFlags.Command, "c", "", "Command: train or predict")
	flag.StringVar(&cmdFlags.InputFile, "i", "", "Path to input CSV file")
	flag.StringVar(&cmdFlags.TargetColumn, "t", "", "Name of target column (only for training)")
	flag.StringVar(&cmdFlags.OutputPath, "o", "", "Path to save output (model or predictions)")
	flag.StringVar(&cmdFlags.TrainedFile, "m", "", "Path to trained model file (only for prediction)")
}

// main parses the flags and prints the values
func main() {
	flag.Parse() // Parse command-line flags
}

func executeTrainCommand(inputPath, targetColumn, outputPath string) error
func executePredictCommand(inputPath, modelPath, outputPath string) error
func parseFlags() (command string, inputPath string, targetColumn string, modelPath string, outputPath string, err error)
