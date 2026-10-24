package gocontainer

import (
	"fmt"
	"maps"
	"reflect"
	"sync"
	"sync/atomic"
)

// Container holds services, each registered under a Key. Put services in and get them out with the key's methods.
//
// The zero value is an empty container, ready to use. A Container is safe for concurrent use
// and must not be copied after first use.
type Container struct {
	mu sync.Mutex // serialises Register and Deregister; reads never take it
	// The registry is copied on write and never modified once stored, so reads need no lock and don't contend:
	// a writer clones it, changes the clone and swaps it in.
	services atomic.Pointer[map[*keyID]any]
}

// load returns the current registry. It is nil for a container nothing was registered in, which reads as empty.
func (c *Container) load() map[*keyID]any {
	if m := c.services.Load(); m != nil {
		return *m
	}
	return nil
}

func (c *Container) set(id *keyID, v any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	next := maps.Clone(c.load())
	if next == nil {
		next = make(map[*keyID]any, 1)
	}
	next[id] = v
	c.services.Store(&next)
}

func (c *Container) remove(id *keyID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	current := c.load()
	if _, ok := current[id]; !ok {
		return
	}
	next := maps.Clone(current)
	delete(next, id)
	c.services.Store(&next)
}

// keyID is what a container finds a service by. Keys are compared by its address, so every NewKey call makes a new key.
type keyID struct{ name string }

// Key names a service of type T. Register takes a T, and Get and MustGet return one, so no type assertions are needed.
//
// Make keys with NewKey. The zero Key names nothing: Register panics with it, and nothing is found under it.
type Key[T any] struct {
	_  [0]*T // makes Key[T] and Key[U] different types, so one can't be converted to the other
	id *keyID
}

// NewKey returns a new key for a service of type T. Every call returns a different key, even for the same name and type,
// so declare each key once, in the package that provides the service, and share it. The name labels the key in panic messages.
func NewKey[T any](name string) Key[T] {
	return Key[T]{id: &keyID{name: name}}
}

// Register puts v in c under k, replacing the service registered under k before, if any.
// It panics if k is the zero Key.
func (k Key[T]) Register(c *Container, v T) {
	if k.id == nil {
		panic("gocontainer: Register with a zero Key; make keys with NewKey")
	}
	c.set(k.id, v)
}

// Deregister removes the service registered under k from c. It does nothing if there is none.
func (k Key[T]) Deregister(c *Container) {
	c.remove(k.id)
}

// Has reports whether a service is registered under k in c.
func (k Key[T]) Has(c *Container) bool {
	_, ok := c.load()[k.id]
	return ok
}

// Get returns the service registered under k in c, and whether there is one.
// When there is none, it returns the zero value of T and false.
func (k Key[T]) Get(c *Container) (T, bool) {
	v, ok := c.load()[k.id]
	// Register stores only a T under k.id, so this fails only for a nil v: no service,
	// or a nil interface that was registered, which comes back nil with ok true.
	t, _ := v.(T)
	return t, ok
}

// MustGet returns the service registered under k in c. It panics if there is none.
func (k Key[T]) MustGet(c *Container) T {
	t, ok := k.Get(c)
	if !ok {
		// An error, formatted only if printed, keeps MustGet small enough to inline.
		panic(notRegistered[T]{k})
	}
	return t
}

// notRegistered is what MustGet panics with.
type notRegistered[T any] struct{ k Key[T] }

func (e notRegistered[T]) Error() string {
	name := ""
	if e.k.id != nil {
		name = e.k.id.name
	}
	return fmt.Sprintf("gocontainer: no %v registered as %q", reflect.TypeFor[T](), name)
}
