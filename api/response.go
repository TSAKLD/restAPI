package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"restAPI/entity"
)

type Error struct {
	Error string `json:"error"`
}

// sendError sending response with error based on error + status code.
func sendError(w http.ResponseWriter, err error) {
	log.Println(err)

	w.Header().Set("Content-Type", "application/json")

	statusCode := http.StatusInternalServerError

	switch {
	case errors.Is(err, entity.ErrNotFound):
		statusCode = http.StatusNotFound
	case errors.Is(err, entity.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
	case errors.Is(err, entity.ErrForbidden):
		statusCode = http.StatusForbidden
	}

	w.WriteHeader(statusCode)

	err = json.NewEncoder(w).Encode(Error{Error: err.Error()})
	if err != nil {
		log.Println(err)
	}
}

func sendResponse(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
