package serialization

import (
	"encoding/json"
	"os"

	"CodeLens-decision-tree/internal/model"
)

// SerializeTree serializes the DecisionTree and writes it to a JSON file.
func SerializeTree(tree *model.DecisionTree, path string) error {
	data, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
