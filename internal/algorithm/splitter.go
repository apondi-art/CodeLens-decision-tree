package algorithm

import (
	"errors"
	"sort"

	"CodeLens-decision-tree/internal/model"
)

// SplitDataset divides the dataset based on a specified attribute value.
// It returns a subset where the attribute matches the provided value.
//
// Parameters:
// - dataset: A pointer to the dataset to be split.
// - attr: A pointer to the attribute used for splitting.
// - value: The attribute value to filter records by.
//
// Returns:
// - A new dataset containing only records matching the given attribute value.
// - An error if dataset or attribute is nil.
