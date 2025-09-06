package response

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	ErrorCode int    `json:"errorcode,omitempty"`
	Data      any    `json:"data,omitempty"`
}

func SuccessResponse(w http.ResponseWriter, data any, message string, code int) {

	response := response{
		Status:  "Success",
		Message: message,
		Data:    data,
	}

	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, errMessage string, code int) {
	response := response{
		Status:    "fail",
		Message:   errMessage,
		ErrorCode: code,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
