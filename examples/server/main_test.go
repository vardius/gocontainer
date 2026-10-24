package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vardius/gocontainer/v2"
)

// downUsers is a fake whose database is unreachable.
type downUsers struct{}

func (downUsers) Name(context.Context, string) (string, error) {
	return "", errors.New("database is down")
}

func TestRoutes(t *testing.T) {
	tests := []struct {
		name     string
		users    Users
		wantCode int
		wantBody string
	}{
		{"found", staticUsers{"1": "Ada"}, http.StatusOK, "Ada\n"},
		{"not found", staticUsers{}, http.StatusNotFound, "404 page not found\n"},
		{"store down", downUsers{}, http.StatusInternalServerError, "database is down\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Each test registers its own Users, in a container of its own.
			var c gocontainer.Container
			usersKey.Register(&c, tt.users)

			w := httptest.NewRecorder()
			routes(&c).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/users/1", nil))

			if w.Code != tt.wantCode || w.Body.String() != tt.wantBody {
				t.Errorf("GET /users/1 = %d %q, want %d %q", w.Code, w.Body, tt.wantCode, tt.wantBody)
			}
		})
	}
}
