package httpx

import (
	"encoding/json"
	"net/http"
)

type Meta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    Meta   `json:"meta,omitzero"`
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	payload any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}

func OK(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func Created(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

func Message(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, Response{
		Success: status < 400,
		Message: msg,
	})
}
