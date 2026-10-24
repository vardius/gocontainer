package gocontainer_test

import (
	"fmt"
	"os"

	"github.com/vardius/gocontainer/v2"
)

// Greeter is the interface the example's services are registered as.
type Greeter interface {
	Greet(name string) string
}

type politeGreeter struct{ greeting string }

func (g politeGreeter) Greet(name string) string { return g.greeting + ", " + name + "." }

func Example() {
	// A key names a service and fixes its type.
	greeterKey := gocontainer.NewKey[Greeter]("greeter")

	// The zero Container is empty and ready to use.
	var c gocontainer.Container
	greeterKey.Register(&c, politeGreeter{greeting: "Good morning"})

	// MustGet returns a Greeter: no type assertion.
	greeter := greeterKey.MustGet(&c)
	fmt.Println(greeter.Greet("Ada"))
	// Output: Good morning, Ada.
}

func ExampleNewKey() {
	// Two keys with the same name and type are still two keys: each NewKey call makes a new one.
	primary := gocontainer.NewKey[string]("dsn")
	replica := gocontainer.NewKey[string]("dsn")

	var c gocontainer.Container
	primary.Register(&c, "postgres://primary")

	fmt.Println(primary.Has(&c), replica.Has(&c))
	// Output: true false
}

func ExampleKey_Register() {
	greeterKey := gocontainer.NewKey[Greeter]("greeter")

	var c gocontainer.Container
	// A politeGreeter goes in as a Greeter: the key's type checks it at compile time.
	greeterKey.Register(&c, politeGreeter{greeting: "Hello"})

	// Registering again under a key replaces the service.
	greeterKey.Register(&c, politeGreeter{greeting: "Hi"})

	fmt.Println(greeterKey.MustGet(&c).Greet("Grace"))
	// Output: Hi, Grace.
}

func ExampleKey_Get() {
	portKey := gocontainer.NewKey[int]("port")

	var c gocontainer.Container
	port, ok := portKey.Get(&c)
	fmt.Println(port, ok)

	portKey.Register(&c, 8080)
	port, ok = portKey.Get(&c)
	fmt.Println(port, ok)
	// Output:
	// 0 false
	// 8080 true
}

func ExampleKey_MustGet() {
	outKey := gocontainer.NewKey[*os.File]("out")

	var c gocontainer.Container
	outKey.Register(&c, os.Stdout)

	fmt.Fprintln(outKey.MustGet(&c), "written to the registered *os.File")
	// Output: written to the registered *os.File
}

func ExampleKey_MustGet_notRegistered() {
	defer func() { fmt.Println(recover()) }()

	var c gocontainer.Container
	gocontainer.NewKey[Greeter]("greeter").MustGet(&c)
	// Output: gocontainer: no gocontainer_test.Greeter registered as "greeter"
}

func ExampleKey_Deregister() {
	tokenKey := gocontainer.NewKey[string]("token")

	var c gocontainer.Container
	tokenKey.Register(&c, "s3cr3t")
	tokenKey.Deregister(&c)

	fmt.Println(tokenKey.Has(&c))
	// Output: false
}
