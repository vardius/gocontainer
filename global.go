package gocontainer

var (
	// InitializeGlobalContainer enables global container instance
	// instance will be initialized with package import by init() function
	InitializeGlobalContainer = true
	c                         Container
)

// Register service by id
func Register(id string, object interface{}) {
	if c == nil {
		panic("Global container instance is not initialized")
	}

	c.Register(id, object)
}

// Deregister service be id
func Deregister(id string) {
	if c == nil {
		panic("Global container instance is not initialized")
	}

	c.Deregister(id)
}

// Has checks if container has an object
func Has(id string) bool {
	if c == nil {
		panic("Global container instance is not initialized")
	}

	return c.Has(id)
}

// Get a service by id
func Get(id string) (interface{}, bool) {
	if c == nil {
		panic("Global container instance is not initialized")
	}

	return c.Get(id)
}

// MustGet calls Get underneath
// will panic if object not found within container
func MustGet(id string) interface{} {
	if c == nil {
		panic("Global container instance is not initialized")
	}

	return c.MustGet(id)
}

// Invoke gets a service safely typed by passing it to a closure
// will panic if callback is not a function
func Invoke(id string, fn interface{}) {
	if c == nil {
		panic("Global container instance is not initialized")
	}

	c.Invoke(id, fn)
}

// MustInvoke calls MustGet underneath
// will panic if object not found within container
func MustInvoke(id string, fn interface{}) {
	if c == nil {
		panic("Global container instance is not initialized")
	}

	c.MustInvoke(id, fn)
}

func init() {
	if !InitializeGlobalContainer {
		return
	}

	c = New()
}
