package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type Error struct {
	Error string `json:"error"`
}

// sendError sending response with error based on error + status code.
func sendError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	err = json.NewEncoder(w).Encode(Error{Error: err.Error()})
	if err != nil {
		log.Println(err)
	}
}

func sendResponse(w http.ResponseWriter, v ...interface{}) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
