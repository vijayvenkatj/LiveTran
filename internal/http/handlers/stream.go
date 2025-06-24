package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (handler *Handler) StopStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	handler.tm.StopTask("1",errors.New("user initiated request"))
	
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    "Stream stopped!",
	})
}

func (handler *Handler) Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")

	task,exists := handler.tm.TaskMap["1"]
	if exists {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(Response{
		Success: true,
		Data: fmt.Sprintf("Status: %s",task.Status),
		})
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error: "Task not found",
	})
}