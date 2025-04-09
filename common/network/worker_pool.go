package network

import (
	"context"
	"fmt"
	"sync"
)

// Task represents a unit of work to be processed
type Task struct {
	Data    []byte
	ConnID  string
	Context context.Context
}

// WorkerPool manages a pool of workers for processing tasks
type WorkerPool struct {
	numWorkers int
	tasks      chan Task
	wg         sync.WaitGroup
	quit       chan struct{}
	handler    func(Task) error
	mu         sync.RWMutex
	closed     bool
}

// NewWorkerPool creates a new worker pool with the specified number of workers
func NewWorkerPool(numWorkers int, bufferSize int, handler func(Task) error) *WorkerPool {
	pool := &WorkerPool{
		numWorkers: numWorkers,
		tasks:      make(chan Task, bufferSize),
		quit:       make(chan struct{}),
		handler:    handler,
	}

	pool.Start()
	return pool
}

// Start launches the worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Submit adds a task to the worker pool
func (wp *WorkerPool) Submit(task Task) error {
	wp.mu.RLock()
	if wp.closed {
		wp.mu.RUnlock()
		return fmt.Errorf("worker pool is closed")
	}
	wp.mu.RUnlock()

	select {
	case wp.tasks <- task:
		return nil
	case <-wp.quit:
		return fmt.Errorf("worker pool is stopping")
	}
}

// worker is the main worker routine
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.quit:
			return
		case task, ok := <-wp.tasks:
			if !ok {
				return
			}
			if err := wp.handler(task); err != nil {
				// Here you could implement error handling, logging, or retry logic
				continue
			}
		}
	}
}

// Stop gracefully shuts down the worker pool
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	if wp.closed {
		wp.mu.Unlock()
		return
	}
	wp.closed = true
	close(wp.quit)
	wp.mu.Unlock()

	wp.wg.Wait()

	wp.mu.Lock()
	close(wp.tasks)
	wp.mu.Unlock()
}

// Len returns the current number of tasks in the queue
func (wp *WorkerPool) Len() int {
	return len(wp.tasks)
}
