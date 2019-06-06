package gocontainer

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	c    *container
	once sync.Once
)

type container struct {
	sync.RWMutex
	registry map[string]interface{}
}

// Register service by id
func Register(id string, object interface{}) {
	if c.registry == nil {
		c.registry = make(map[string]interface{})
	}
	c.Lock()
	c.registry[id] = object
	c.Unlock()
}

// Deregister service be id
func Deregister(id string) {
	c.Lock()
	delete(c.registry, id)
	c.Unlock()
}

// Has checks if container has an object
func Has(id string) bool {
	_, ok := Get(id)

	return ok
}

// Get a service by id
func Get(id string) (interface{}, bool) {
	c.RLock()
	o, ok := c.registry[id]
	c.RUnlock()

	return o, ok
}

// MustGet calls Get underneath
// will panic if object not found within container
func MustGet(id string) interface{} {
	o, ok := Get(id)
	if !ok {
		panic(fmt.Sprintf("Object <%s> nof found within a container", id))
	}

	return o
}

// Invoke gets a service safely typed by passing it to a closure
// will panic if callback is not a function
func Invoke(id string, fn interface{}) {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		panic(fmt.Sprintf("%s is not a reflect.Func", reflect.TypeOf(fn)))
	}

	callback := reflect.ValueOf(fn)
	o, ok := Get(id)

	args := []reflect.Value{reflect.ValueOf(o), reflect.ValueOf(ok)}

	callback.Call(args)
}

// MustInvoke calls Invoke underneath
// will panic if object not found within container
func MustInvoke(id string, fn interface{}) {
	if !Has(id) {
		panic(fmt.Sprintf("Object <%s> nof found within a container", id))
	}

	Invoke(id, fn)
}

func init() {
	once.Do(func() {
		c = &container{
			registry: make(map[string]interface{}),
		}
	})
}
