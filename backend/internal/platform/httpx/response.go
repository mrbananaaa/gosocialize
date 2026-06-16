package httpx

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitzero"`
}

type PaginationMeta struct {
	NextCursor string `json:"next_cursor"` // base64encoded
	HasMore    bool   `json:"has_more"`
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

func OK(w http.ResponseWriter, data any, opts ...ResponseOptions) {
	res := &Response{
		Success: true,
		Data:    data,
	}

	for _, opt := range opts {
		opt(res)
	}

	writeJSON(w, http.StatusOK, res)
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

type ResponseOptions func(*Response)

func WithMeta(m any) ResponseOptions {
	return func(r *Response) {
		r.Meta = m
	}
}
