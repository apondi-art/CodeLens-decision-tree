package parallel

import "CodeLens-decision-tree/internal/model"

type (
	Task       func() interface{}
	WorkerPool struct {
		// Fields to manage worker pool
	}
)

func NewWorkerPool(numWorkers int) *WorkerPool
func (wp *WorkerPool) Submit(task Task) interface{}
func (wp *WorkerPool) ProcessInParallel(tasks []Task) []interface{}

// Parallel tree building
func BuildTreeParallel(dataset *model.Dataset, attributes []*model.Attribute, targetAttr string, depth, maxDepth int, workerPool *WorkerPool) (*model.Node, error)
