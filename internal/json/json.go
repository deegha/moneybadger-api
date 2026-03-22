package json

import (
	"encoding/json"
	"net/http"
)

type Data struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func Writer(w http.ResponseWriter, status int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := &Data{
		Data:    data,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

func Reader(r *http.Request, data interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}
