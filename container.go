package gocontainer

import (
	"fmt"
	"reflect"
	"sync"
)

// Container interface
type Container interface {
	Register(id string, object interface{})
	Deregister(id string)
	Has(id string) bool
	Get(id string) (interface{}, bool)
	MustGet(id string) interface{}
	Invoke(id string, fn interface{})
	MustInvoke(id string, fn interface{})
	MustInvokeMany(ids ...string) func(fn interface{})
}

type container struct {
	sync.RWMutex
	registry map[string]interface{}
}

// New creates new container
func New() Container {
	return &container{
		registry: make(map[string]interface{}),
	}
}

// Register service by id
func (c *container) Register(id string, object interface{}) {
	if c.registry == nil {
		c.registry = make(map[string]interface{})
	}
	c.Lock()
	c.registry[id] = object
	c.Unlock()
}

// Deregister service be id
func (c *container) Deregister(id string) {
	c.Lock()
	delete(c.registry, id)
	c.Unlock()
}

// Has checks if container has an object
func (c *container) Has(id string) bool {
	_, ok := c.Get(id)

	return ok
}

// Get a service by id
func (c *container) Get(id string) (interface{}, bool) {
	c.RLock()
	o, ok := c.registry[id]
	c.RUnlock()

	return o, ok
}

// MustGet calls Get underneath
// will panic if object not found within container
func (c *container) MustGet(id string) interface{} {
	o, ok := c.Get(id)
	if !ok {
		panic(fmt.Sprintf("Object <%s> nof found within a container", id))
	}

	return o
}

// Invoke gets a service safely typed by passing it to a closure
// will panic if callback is not a function
func (c *container) Invoke(id string, fn interface{}) {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		panic(fmt.Sprintf("%s is not a reflect.Func", reflect.TypeOf(fn)))
	}

	o, ok := c.Get(id)
	callback := reflect.ValueOf(fn)
	args := []reflect.Value{reflect.ValueOf(o), reflect.ValueOf(ok)}

	callback.Call(args)
}

// MustInvoke calls MustGet underneath
// will panic if object not found within container
func (c *container) MustInvoke(id string, fn interface{}) {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		panic(fmt.Sprintf("%s is not a reflect.Func", reflect.TypeOf(fn)))
	}

	o := c.MustGet(id)

	callback := reflect.ValueOf(fn)
	args := []reflect.Value{reflect.ValueOf(o)}

	callback.Call(args)
}

// MustInvokeMany calls MustInvoke underneath
// returns many services from container
// will panic if object not found within container
func (c *container) MustInvokeMany(ids ...string) func(fn interface{}) {
	var args []reflect.Value

	for _, id := range ids {
		o := c.MustGet(id)

		args = append(args, reflect.ValueOf(o))
	}

	return func(fn interface{}) {
		if reflect.TypeOf(fn).Kind() != reflect.Func {
			panic(fmt.Sprintf("%s is not a reflect.Func", reflect.TypeOf(fn)))
		}

		callback := reflect.ValueOf(fn)
		callback.Call(args)
	}
}
