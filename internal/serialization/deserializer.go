package serialization

import "CodeLens-decision-tree/internal/model"

func DeserializeTree(path string) (*model.DecisionTree, error)
func JSONToNode(jsonData map[string]interface{}) (*model.Node, error)
