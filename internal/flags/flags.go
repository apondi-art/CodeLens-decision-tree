package flags

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Flags stores command-line arguments
type Flags struct {
	Command      string // train or predict
	InputFile    string // Path to input CSV file
	TargetColumn string // Only for training
	OutputPath   string // Model output or predictions
	TrainedFile  string // Only for prediction
}


// ParseFlags initializes and validates CLI flags
func ParseFlags() (Flags, error) {
	var cmdFlags Flags

	// Create a new flag set (instead of using the default `flag` package)
	flagSet := flag.NewFlagSet("dt", flag.ContinueOnError)

	flagSet.StringVar(&cmdFlags.Command, "c", "", "Command: train or predict")
	flagSet.StringVar(&cmdFlags.InputFile, "i", "", "Path to input CSV file")
	flagSet.StringVar(&cmdFlags.TargetColumn, "t", "", "Target column (only for training)")
	flagSet.StringVar(&cmdFlags.OutputPath, "o", "", "Output file path (model or predictions)")
	flagSet.StringVar(&cmdFlags.TrainedFile, "m", "", "Path to trained model file (only for prediction)")

	// Capture unexpected flag errors
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		return Flags{}, fmt.Errorf("invalid flag provided: %w", err)
	}

	// Validate flags before returning
	if err := validateFlags(cmdFlags); err != nil {
		return Flags{}, err
	}

	return cmdFlags, nil
}

// validateFlags ensures all required flags are provided with proper error messages
func validateFlags(cmdFlags Flags) error {
	var missingFlags []string

	// Command must be provided
	if cmdFlags.Command == "" {
		return errors.New("missing required flag: -c (train or predict). Example: -c train")
	}

	switch cmdFlags.Command {
	case "train":
		missingFlags = collectMissingFlags(cmdFlags.InputFile, "-i <input.csv>", missingFlags)
		missingFlags = collectMissingFlags(cmdFlags.TargetColumn, "-t <target_column>", missingFlags)
		missingFlags = collectMissingFlags(cmdFlags.OutputPath, "-o <output.dt>", missingFlags)

	case "predict":
		missingFlags = collectMissingFlags(cmdFlags.InputFile, "-i <test.csv>", missingFlags)
		missingFlags = collectMissingFlags(cmdFlags.TrainedFile, "-m <model.dt>", missingFlags)
		missingFlags = collectMissingFlags(cmdFlags.OutputPath, "-o <predictions.csv>", missingFlags)

	default:
		return fmt.Errorf("invalid command: %q. Use -c train or -c predict", cmdFlags.Command)
	}

	// If any flags are missing, returns a single error message
	if len(missingFlags) > 0 {
		return fmt.Errorf(
			"missing required flags for %s: %s\nExample usage:\n  Train: ./dt -c train -i data.csv -t class -o model.dt\n  Predict: ./dt -c predict -i test.csv -m model.dt -o predictions.csv",
			cmdFlags.Command, strings.Join(missingFlags, ", "),
		)
	}

	return nil
}

// collectMissingFlags is a helper function to check missing flags
func collectMissingFlags(value, flagName string, missingFlags []string) []string {
	if value == "" {
		missingFlags = append(missingFlags, flagName)
	}
	return missingFlags
}
