package main

import (
	"fmt"
	api "github.com/vijayvenkatj/LiveTran/internal/http"
	"github.com/vijayvenkatj/LiveTran/internal/ingest"
)

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