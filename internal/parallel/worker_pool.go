package parallel

import "sync"

type (
	Task func() interface{}

	WorkerPool struct {
		numWorkers  int
		jobQueue    chan Task
		resultQueue chan interface{}
		wg          sync.WaitGroup
	}
)

func NewWorkerPool(numWorkers int) *WorkerPool {
	wp := &WorkerPool{
		numWorkers:  numWorkers,
		jobQueue:    make(chan Task, 100),        // Buffered channel
		resultQueue: make(chan interface{}, 100), // Buffered channel
		wg:          sync.WaitGroup{},
	}

	for i := 0; i < numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	return wp
}

func (wp *WorkerPool) Submit(task Task) {
	wp.jobQueue <- task
}

func (wp *WorkerPool) ProcessInParallel(tasks []Task) []interface{} {
	results := make([]interface{}, len(tasks))
	resultChan := make(chan interface{}, len(tasks))

	// Submit tasks to the job queue
	for _, task := range tasks {
		wp.Submit(func() interface{} {
			return task()
		})
	}

	// Close the job queue after all tasks have been submitted
	close(wp.jobQueue)

	// Collect results from the result queue
	go func() {
		for result := range wp.resultQueue {
			resultChan <- result
		}
		close(resultChan)
	}()

	// Wait for all workers to complete
	wp.wg.Wait()
	close(wp.resultQueue)

	// Collect results from the result channel
	i := 0
	for result := range resultChan {
		results[i] = result
		i++
	}

	return results
}

func (wp *WorkerPool) worker(workerID int) {
	defer wp.wg.Done()

	for task := range wp.jobQueue {
		result := task()
		wp.resultQueue <- result
	}
}

// func NewWorkerPool(numWorkers int) *WorkerPool
// func (wp *WorkerPool) Submit(task Task) interface{}
// func (wp *WorkerPool) ProcessInParallel(tasks []Task) []interface{}

// Parallel tree building
// func BuildTreeParallel(dataset *model.Dataset, attributes []*Attribute, targetAttr string, depth, maxDepth int, workerPool *WorkerPool) (*model.Node, error)
