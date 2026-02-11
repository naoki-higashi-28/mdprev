package handler

import (
	"errors"
	"net/http"

	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrForbidden):
		http.Error(w, "Forbidden", http.StatusForbidden)
	case errors.Is(err, model.ErrNotFound):
		http.Error(w, "Not Found", http.StatusNotFound)
	case errors.Is(err, model.ErrBadRequest), errors.Is(err, model.ErrInvalidExtFile):
		http.Error(w, "Bad Request", http.StatusBadRequest)
	default:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
