package gocontainer_test

import (
	"fmt"
	"testing"

	"github.com/vardius/gocontainer/v2"
)

// fill returns a container with n *service entries and their keys.
func fill(n int) (*gocontainer.Container, []gocontainer.Key[*service]) {
	c := new(gocontainer.Container)
	keys := make([]gocontainer.Key[*service], n)
	for i := range keys {
		keys[i] = gocontainer.NewKey[*service](fmt.Sprint("service ", i))
		keys[i].Register(c, &service{name: fmt.Sprint(i)})
	}
	return c, keys
}

// TestAllocations guards the allocations per call. Lower a number when an optimisation lands;
// never raise one without a reason in the commit.
func TestAllocations(t *testing.T) {
	if raceEnabled {
		t.Skip("the race detector changes allocation counts")
	}

	c, keys := fill(32)
	k := keys[16]
	missing := gocontainer.NewKey[*service]("missing")
	empty := new(gocontainer.Container)
	one, oneKeys := fill(1)
	port := gocontainer.NewKey[int]("port")
	port.Register(one, 0)
	db, n := &service{}, 1000

	tests := []struct {
		name string
		call func()
		want float64
	}{
		{"Get", func() { _, _ = k.Get(c) }, 0},
		{"Get, not registered", func() { _, _ = missing.Get(c) }, 0},
		{"Get, empty container", func() { _, _ = k.Get(empty) }, 0},
		{"MustGet", func() { _ = k.MustGet(c) }, 0},
		{"Has", func() { _ = k.Has(c) }, 0},
		{"Deregister, not registered", func() { missing.Deregister(c) }, 0},
		// The new registry, small enough for one group of slots, and the pointer to it that is swapped in.
		{"Register a pointer, 2 services", func() { oneKeys[0].Register(one, db) }, 3},
		// Plus the int, boxed into an interface: values that aren't pointers are copied to the heap.
		{"Register an int, 2 services", func() { n++; port.Register(one, n) }, 4},
		// A registry this size copies as four allocations, its directory, table and slots among them, plus the pointer.
		{"Register a pointer, 32 services", func() { k.Register(c, db) }, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testing.AllocsPerRun(1000, tt.call); got != tt.want {
				t.Errorf("allocates %v times per call, want %v", got, tt.want)
			}
		})
	}
}
