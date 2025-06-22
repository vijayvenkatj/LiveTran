package ingest

import (
	"context"
	"fmt"
	"sync"
)


type TaskManager struct {
	mu		sync.Mutex
	wg 	    *sync.WaitGroup
	TaskMap	map[string]context.CancelCauseFunc
}


// Constuctor

func NewTaskManager(wg *sync.WaitGroup) *TaskManager {
	return &TaskManager{
		wg: wg,
		TaskMap: make(map[string]context.CancelCauseFunc),
	}
}


// Starting a Task 

func (tm *TaskManager) StartTask(id string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, exists := tm.TaskMap[id]; exists {
		fmt.Println("Job exists!")
		// Re-Connection logic ???
		return 
	}

	cancelCtx,cancelFunc := context.WithCancelCause(context.Background())
	tm.TaskMap[id] = cancelFunc

	tm.wg.Add(1)

	go func() {
		defer tm.wg.Done()
		SrtConnectionTask(cancelCtx,id)
	}()
}

// Stopping a task

func (tm *TaskManager) StopTask(id string,reason error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if cancelFunc, exists := tm.TaskMap[id]; exists {
		cancelFunc(reason)
		delete(tm.TaskMap,id)
	} else {
		fmt.Println("Job already done / Cancelled")
	}

}