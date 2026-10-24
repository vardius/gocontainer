// Command server serves users' names over HTTP. Its handler takes the user store from a container,
// so a test can hand it a fake.
//
//	go run ./examples/server
//	curl localhost:8080/users/1
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/vardius/gocontainer/v2"
)

// ErrNotFound is what Users returns for an id it doesn't know.
var ErrNotFound = errors.New("user not found")

// Users looks up a user's name by id.
type Users interface {
	Name(ctx context.Context, id string) (string, error)
}

var usersKey = gocontainer.NewKey[Users]("users")

// routes returns the app's handler, built from the services in c.
func routes(c *gocontainer.Container) http.Handler {
	users := usersKey.MustGet(c)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		name, err := users.Name(r.Context(), r.PathValue("id"))
		if errors.Is(err, ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, name)
	})
	return mux
}

// staticUsers is the Users the server runs with: a fixed list.
type staticUsers map[string]string

func (s staticUsers) Name(_ context.Context, id string) (string, error) {
	name, ok := s[id]
	if !ok {
		return "", ErrNotFound
	}
	return name, nil
}

func main() {
	var c gocontainer.Container
	usersKey.Register(&c, staticUsers{"1": "Ada", "2": "Grace"})

	log.Fatal(http.ListenAndServe("localhost:8080", routes(&c)))
}
