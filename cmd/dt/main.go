package main

import (
	"errors"
	"flag"
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
