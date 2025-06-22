package main

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"github.com/vijayvenkatj/LiveTran/internal/ingest"
)

func main() {

	var wg sync.WaitGroup

	tm := ingest.NewTaskManager(&wg)

	tm.StartTask("1");
	tm.StartTask("2")
	
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	tm.StopTask("1",errors.New("ShutDown Gracefully!"));
	tm.StopTask("2",errors.New("ShutDown Gracefully!"));
	wg.Wait()
}