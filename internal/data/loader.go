package data

import "CodeLens-decision-tree/internal/model"

func LoadCSV(path string) (*model.Dataset, error)
func InferDataTypes(data [][]string) (map[string]string, error)
func HandleMissingValues(dataset *model.Dataset) error
