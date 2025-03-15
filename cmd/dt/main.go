package main

import (
	"fmt"
	"log"

	"CodeLens-decision-tree/internal/flags"
)

// This is the entry point of the program
func main() {
	// Parse CLI flags.
	cmdFlags, err := flags.ParseFlags()
	if err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	// Execute the appropriate command.
	switch cmdFlags.Command {
	case "train":
		err = flags.ExecuteTrainCommand(cmdFlags.InputFilePath, cmdFlags.TargetColumn, cmdFlags.OutputPath)
	case "predict":
		err = flags.ExecutePredictCommand(cmdFlags.InputFilePath, cmdFlags.TrainedFilePath, cmdFlags.OutputPath)
	default:
		log.Fatalf("Invalid command: %s", cmdFlags.Command)
	}

	if err != nil {
		log.Fatalf("Error executing command: %v", err)
	}

	fmt.Println("Command executed successfully.")
}
