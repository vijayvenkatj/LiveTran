package ingest

import (
	"context"
	"fmt"
	"sync"
)


type TaskManager struct {
	mu		sync.Mutex
	TaskMap	map[string]context.CancelCauseFunc
}


// Constuctor

func NewTaskManager() *TaskManager {
	return &TaskManager{
		TaskMap: make(map[string]context.CancelCauseFunc),
	}
}


// Starting a Task 

func (tm *TaskManager) StartTask(id string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, exists := tm.TaskMap[id]; exists {
		fmt.Println("Job exists!")
		return 
	}

	cancelCtx,cancelFunc := context.WithCancelCause(context.Background())
	tm.TaskMap[id] = cancelFunc


	go func() {
		err := SrtConnectionTask(cancelCtx,id)
		if err != nil {
			fmt.Println(err)
			// Reconnection logic??
		}
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