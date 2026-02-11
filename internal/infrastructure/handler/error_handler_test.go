package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
)

func TestWriteError(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
	}{
		{name: "forbidden", err: model.ErrForbidden, status: http.StatusForbidden},
		{name: "not found", err: model.ErrNotFound, status: http.StatusNotFound},
		{name: "bad request", err: model.ErrBadRequest, status: http.StatusBadRequest},
		{name: "invalid extension", err: model.ErrInvalidExtFile, status: http.StatusBadRequest},
		{name: "internal", err: errors.New("boom"), status: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			writeError(rr, tt.err)
			if rr.Code != tt.status {
				t.Fatalf("status = %d, want %d", rr.Code, tt.status)
			}
		})
	}
}
