package main

import (
	"fmt"
	api "github.com/vijayvenkatj/LiveTran/internal/http"
	"github.com/vijayvenkatj/LiveTran/internal/ingest"
)

// func main() {

// 	var wg sync.WaitGroup

// 	tm := ingest.NewTaskManager(&wg)

// 	tm.StartTask("1");
// 	tm.StartTask("2")

// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
// 	<-c

// 	tm.StopTask("1",errors.New("ShutDown Gracefully!"));
// 	tm.StopTask("2",errors.New("ShutDown Gracefully!"));
// 	wg.Wait()
// }

var tm *ingest.TaskManager

func init() {
	tm = ingest.NewTaskManager()
}

func main() {
	apiServer := api.NewAPIServer(":8080")
	err := apiServer.StartAPIServer(tm);
	if err != nil {
		fmt.Println("Error starting server")
		return
	}
}