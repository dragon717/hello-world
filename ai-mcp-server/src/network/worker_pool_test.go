package network

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewWorkerPool(t *testing.T) {
	handler := func(task Task) error {
		return nil
	}

	tests := []struct {
		name       string
		numWorkers int
		bufferSize int
		handler    func(Task) error
	}{
		{
			name:       "Valid configuration",
			numWorkers: 5,
			bufferSize: 10,
			handler:    handler,
		},
		{
			name:       "Minimum configuration",
			numWorkers: 1,
			bufferSize: 1,
			handler:    handler,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := NewWorkerPool(tt.numWorkers, tt.bufferSize, tt.handler)
			if pool == nil {
				t.Error("NewWorkerPool() returned nil pool")
			}
			if pool.numWorkers != tt.numWorkers {
				t.Errorf("NewWorkerPool() numWorkers = %v, want %v", pool.numWorkers, tt.numWorkers)
			}
			if cap(pool.tasks) != tt.bufferSize {
				t.Errorf("NewWorkerPool() buffer size = %v, want %v", cap(pool.tasks), tt.bufferSize)
			}
		})
	}
}

func TestWorkerPool_Submit(t *testing.T) {
	var processedTasks int
	var mu sync.Mutex

	handler := func(task Task) error {
		mu.Lock()
		processedTasks++
		mu.Unlock()
		return nil
	}

	pool := NewWorkerPool(4, 10, handler)

	// Submit tasks
	numTasks := 10
	for i := 0; i < numTasks; i++ {
		err := pool.Submit(Task{
			Data:    []byte("test data"),
			ConnID:  "test",
			Context: context.Background(),
		})
		if err != nil {
			t.Errorf("Failed to submit task: %v", err)
		}
	}

	// Wait for tasks to be processed
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	if processedTasks != numTasks {
		t.Errorf("Expected %d tasks to be processed, got %d", numTasks, processedTasks)
	}
	mu.Unlock()
}

func TestWorkerPool_Stop(t *testing.T) {
	var processedTasks int
	var mu sync.Mutex

	handler := func(task Task) error {
		time.Sleep(10 * time.Millisecond) // Simulate work
		mu.Lock()
		processedTasks++
		mu.Unlock()
		return nil
	}

	pool := NewWorkerPool(2, 5, handler)

	// Submit some tasks
	numTasks := 5
	for i := 0; i < numTasks; i++ {
		err := pool.Submit(Task{
			Data:    []byte("test data"),
			ConnID:  "test",
			Context: context.Background(),
		})
		if err != nil {
			t.Errorf("Failed to submit task: %v", err)
		}
	}

	// Stop the pool
	pool.Stop()

	// Try to submit after stopping
	err := pool.Submit(Task{
		Data:    []byte("test data"),
		ConnID:  "test",
		Context: context.Background(),
	})
	if err == nil {
		t.Error("Expected error when submitting to stopped pool")
	}

	mu.Lock()
	if processedTasks > numTasks {
		t.Errorf("Expected at most %d tasks to be processed, got %d", numTasks, processedTasks)
	}
	mu.Unlock()
}

func TestWorkerPool_Concurrent(t *testing.T) {
	var processedTasks int64
	handler := func(task Task) error {
		time.Sleep(time.Millisecond) // Simulate work
		atomic.AddInt64(&processedTasks, 1)
		return nil
	}

	pool := NewWorkerPool(8, 20, handler)

	// Submit tasks concurrently
	var wg sync.WaitGroup
	numGoroutines := 10
	tasksPerGoroutine := 10
	errorChan := make(chan error, numGoroutines*tasksPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < tasksPerGoroutine; j++ {
				err := pool.Submit(Task{
					Data:    []byte("test data"),
					ConnID:  "test",
					Context: context.Background(),
				})
				if err != nil {
					errorChan <- err
				}
			}
		}()
	}

	// Wait for all submissions to complete
	wg.Wait()
	close(errorChan)

	// Check for any submission errors
	for err := range errorChan {
		t.Errorf("Task submission error: %v", err)
	}

	// Wait for all tasks to be processed
	time.Sleep(500 * time.Millisecond)

	// Verify the number of processed tasks
	totalTasks := int64(numGoroutines * tasksPerGoroutine)
	if atomic.LoadInt64(&processedTasks) != totalTasks {
		t.Errorf("Expected %d tasks to be processed, got %d", totalTasks, atomic.LoadInt64(&processedTasks))
	}

	// Verify the length of the task queue
	queueLen := pool.Len()
	if queueLen > cap(pool.tasks) {
		t.Errorf("Task queue length %d exceeds capacity %d", queueLen, cap(pool.tasks))
	}
}
