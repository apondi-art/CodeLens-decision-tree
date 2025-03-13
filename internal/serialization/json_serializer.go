package serialization

import "CodeLens-decision-tree/internal/model"

func SerializeTree(tree *model.DecisionTree, path string) error
func NodeToJSON(node *model.Node) (map[string]interface{}, error)
