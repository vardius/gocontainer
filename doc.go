// Package gocontainer is a dependency injection container with typed keys.
//
// A Key names a service and fixes its type. Register a service in a Container under its key,
// and get it back as that type, with no type assertion:
//
//	var Users = gocontainer.NewKey[UserStore]("users")
//
//	var c gocontainer.Container
//	Users.Register(&c, NewSQLUserStore(db))
//
//	users := Users.MustGet(&c) // a UserStore
//
// # Keys
//
// NewKey makes a key for a service of a type. Register takes a value of that type, or one assignable to it,
// such as a concrete type for an interface key; anything else is a compile error. Get and MustGet return the type.
//
// Each NewKey call makes a different key, even for the same name and type, as errors.New makes a different error,
// so two packages can't take each other's service by choosing the same name. Declare each key once, usually in the
// package that provides the service, and share it. The name only labels the key in panic messages.
//
// # Containers
//
// The zero Container is empty and ready to use. Declare one where the services are wired, in main or a package of
// its own, and register each service before the code that needs it runs. Register replaces whatever was registered
// under the key, Deregister removes it, and Has reports whether there is one. MustGet panics when there is none;
// Get reports it instead.
//
// A Container is safe for concurrent use. A key can be used with any number of containers, which is how a test
// swaps a service for a fake: register the fake under the same key in a container of its own.
//
// # Cost
//
// Get, MustGet and Has allocate nothing and take no lock: a map lookup and, to return a T, a type assertion.
// Goroutines reading at once don't slow each other down.
//
// Register and Deregister change a copy of the registry and swap it in, so that reads never wait for them.
// They take time in proportion to the services already registered, and allocate the copy: they suit services
// registered at startup, not a registry that changes on every request. A value that isn't a pointer, such as a struct,
// may also be boxed on the heap when it is registered; Get copies it out without allocating.
package gocontainer
