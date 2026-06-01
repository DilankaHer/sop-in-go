package utility

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   any    `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, message string, payload any) {
	if status == 0 {
		status = http.StatusOK
	}
	if message == "" {
		message = http.StatusText(status)
	}

	resp := Response{
		Status:  status,
		Message: message,
		Data:    nil,
		Error:   nil,
	}

	if status >= http.StatusBadRequest {
		resp.Error = payload
	} else {
		resp.Data = payload
	}

	body, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
