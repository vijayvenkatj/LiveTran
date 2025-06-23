package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Data    string `json:"data,omitempty"`
}

func (handler *Handler) StartStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	handler.tm.StartTask("1")
	
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    "Stream launching!",
	})
}

func (handler *Handler) StopStream (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	handler.tm.StopTask("1",errors.New("user initiated request"))
	
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    "Stream stopped!",
	})
}