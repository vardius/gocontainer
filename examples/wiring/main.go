// Command wiring registers an app's services at startup, each built from the ones registered before it.
//
//	go run ./examples/wiring
package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/vardius/gocontainer/v2"
)

// Users looks up a user's name by id.
type Users interface {
	Name(id int) (string, bool)
}

// memoryUsers is a Users kept in a map.
type memoryUsers map[int]string

func (m memoryUsers) Name(id int) (string, bool) {
	name, ok := m[id]
	return name, ok
}

// Greeter greets users by name, and logs the ids it doesn't know.
type Greeter struct {
	users Users
	log   *slog.Logger
}

func (g *Greeter) Greet(id int) string {
	name, ok := g.users.Name(id)
	if !ok {
		g.log.Warn("unknown user", "id", id)
		return "Hello, stranger."
	}
	return "Hello, " + name + "."
}

// Each key names a service and fixes its type.
var (
	logKey     = gocontainer.NewKey[*slog.Logger]("log")
	usersKey   = gocontainer.NewKey[Users]("users")
	greeterKey = gocontainer.NewKey[*Greeter]("greeter")
)

// register adds the app's services to c, building each from those before it.
func register(c *gocontainer.Container, logs io.Writer) {
	logKey.Register(c, slog.New(slog.NewTextHandler(logs, nil)))
	// A memoryUsers goes in under the Users key: no conversion needed.
	usersKey.Register(c, memoryUsers{1: "Ada", 2: "Grace"})
	greeterKey.Register(c, &Greeter{
		users: usersKey.MustGet(c),
		log:   logKey.MustGet(c),
	})
}

func main() {
	var c gocontainer.Container
	register(&c, os.Stderr)

	greeter := greeterKey.MustGet(&c)
	for _, id := range []int{1, 2, 3} {
		fmt.Println(greeter.Greet(id))
	}
}
