package render

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, body any, statusCode int) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	if err := enc.Encode(body); err != nil {
		JSONError(w, "json encode error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write(buf.Bytes())
}

func JSONError(w http.ResponseWriter, msg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}
