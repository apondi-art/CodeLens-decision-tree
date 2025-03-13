package main

// cmd/dt/main.go
func main()
func executeTrainCommand(inputPath, targetColumn, outputPath string) error
func executePredictCommand(inputPath, modelPath, outputPath string) error
func parseFlags() (command string, inputPath string, targetColumn string, modelPath string, outputPath string, err error)
