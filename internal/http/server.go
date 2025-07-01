package api

import (
	"fmt"
	"net/http"

	"github.com/vijayvenkatj/LiveTran/internal/http/handlers"
	"github.com/vijayvenkatj/LiveTran/internal/ingest"
)



type APIServer struct {
	address string
}


// Constructor for APIServer
func NewAPIServer(address string) *APIServer {
	return &APIServer{
		address: address,
	}
}


func (a *APIServer) StartAPIServer(tm *ingest.TaskManager) error {

	routeHandler := handlers.NewHandler(tm)

	streamRoutes := routeHandler.StreamRoutes()

	router := http.NewServeMux()
	router.Handle("/api/",http.StripPrefix("/api",streamRoutes))

	fmt.Println("Server is listening on port",a.address)
	return http.ListenAndServe(a.address,router)

}