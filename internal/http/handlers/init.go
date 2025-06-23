package handlers

import (
	"net/http"
	"github.com/vijayvenkatj/LiveTran/internal/ingest"
)


type Handler struct {
	tm		*ingest.TaskManager
}


// Constructor for Handler
func NewHandler(tm *ingest.TaskManager) *Handler {
	return &Handler{
		tm: tm,
	}
}

// Sub-Router for Stream APIs
func (h *Handler) StreamRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /start-stream",h.StartStream)
	mux.HandleFunc("GET /stop-stream",h.StopStream)


	return mux
}

/*
	This File manages all the SubRouters 
	Flow : 
		APIServer => RouteHandler => Individual Route Handlers based on Path
*/

