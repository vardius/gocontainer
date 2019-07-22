package gocontainer

var (
	// GlobalContainer global container instance
	// instance initialized by package's init() function
	GlobalContainer Container
)

// Register service by id
func Register(id string, object interface{}) {
	if GlobalContainer == nil {
		panic("Global container instance is not initialized")
	}

	GlobalContainer.Register(id, object)
}

// Deregister service be id
func Deregister(id string) {
	if GlobalContainer == nil {
		panic("Global container instance is not initialized")
	}

	GlobalContainer.Deregister(id)
}

// Has checks if container has an object
func Has(id string) bool {
	if GlobalContainer == nil {
		panic("Global container instance is not initialized")
	}

	return GlobalContainer.Has(id)
}

// Get a service by id
func Get(id string) (interface{}, bool) {
	if GlobalContainer == nil {
		panic("Global container instance is not initialized")
	}

	return GlobalContainer.Get(id)
}

// MustGet calls Get underneath
// will panic if object not found within container
func MustGet(id string) interface{} {
	if GlobalContainer == nil {
		panic("Global container instance is not initialized")
	}

	return GlobalContainer.MustGet(id)
}

// Invoke gets a service safely typed by passing it to a closure
// will panic if callback is not a function
func Invoke(id string, fn interface{}) {
	if GlobalContainer == nil {
		panic("Global container instance is not initialized")
	}

	GlobalContainer.Invoke(id, fn)
}

// MustInvoke calls MustGet underneath
// will panic if object not found within container
func MustInvoke(id string, fn interface{}) {
	if GlobalContainer == nil {
		panic("Global container instance is not initialized")
	}

	GlobalContainer.MustInvoke(id, fn)
}

// MustInvokeMany calls MustInvoke underneath
// returns many services from container
// will panic if object not found within container
func MustInvokeMany(ids ...string) func(fn interface{}) {
	return GlobalContainer.MustInvokeMany(ids...)
}

func init() {
	GlobalContainer = New()
}
