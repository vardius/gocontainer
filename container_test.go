package gocontainer_test

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/vardius/gocontainer/v2"
)

type service struct{ name string }

// wantPanic calls f and returns the message it panicked with, failing t if it didn't panic.
func wantPanic(t *testing.T, f func()) (msg string) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Error("did not panic")
		}
		msg = fmt.Sprint(r)
	}()
	f()
	return ""
}

// The key's type parameter is what makes Get typed, so a key must not convert to a key of another type.
func TestKeysOfDifferentTypesDontConvert(t *testing.T) {
	if reflect.TypeFor[gocontainer.Key[int]]().ConvertibleTo(reflect.TypeFor[gocontainer.Key[string]]()) {
		t.Error("a Key[int] converts to a Key[string]")
	}
}

func TestGet(t *testing.T) {
	db := &service{name: "db"}

	tests := []struct {
		name   string
		setup  func(c *gocontainer.Container, k gocontainer.Key[*service])
		want   *service
		wantOK bool
	}{
		{
			name:   "registered",
			setup:  func(c *gocontainer.Container, k gocontainer.Key[*service]) { k.Register(c, db) },
			want:   db,
			wantOK: true,
		},
		{
			// v1's Invoke panicked inside reflect here, instead of reporting the service missing.
			name:  "not registered",
			setup: func(*gocontainer.Container, gocontainer.Key[*service]) {},
		},
		{
			name:   "nil registered",
			setup:  func(c *gocontainer.Container, k gocontainer.Key[*service]) { k.Register(c, nil) },
			wantOK: true,
		},
		{
			name: "deregistered",
			setup: func(c *gocontainer.Container, k gocontainer.Key[*service]) {
				k.Register(c, db)
				k.Deregister(c)
			},
		},
		{
			name: "registered again after deregistering",
			setup: func(c *gocontainer.Container, k gocontainer.Key[*service]) {
				k.Register(c, &service{name: "old"})
				k.Deregister(c)
				k.Register(c, db)
			},
			want:   db,
			wantOK: true,
		},
		{
			name: "replaced",
			setup: func(c *gocontainer.Container, k gocontainer.Key[*service]) {
				k.Register(c, &service{name: "old"})
				k.Register(c, db)
			},
			want:   db,
			wantOK: true,
		},
		{
			name: "only another key with the same name registered",
			setup: func(c *gocontainer.Container, _ gocontainer.Key[*service]) {
				gocontainer.NewKey[*service]("db").Register(c, db)
			},
		},
		{
			name: "only a key of another type with the same name registered",
			setup: func(c *gocontainer.Container, _ gocontainer.Key[*service]) {
				gocontainer.NewKey[string]("db").Register(c, "db")
			},
		},
		{
			name: "only in another container",
			setup: func(_ *gocontainer.Container, k gocontainer.Key[*service]) {
				k.Register(new(gocontainer.Container), db)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c gocontainer.Container
			k := gocontainer.NewKey[*service]("db")
			tt.setup(&c, k)

			got, ok := k.Get(&c)
			if got != tt.want || ok != tt.wantOK {
				t.Errorf("Get() = %v, %v, want %v, %v", got, ok, tt.want, tt.wantOK)
			}
			if has := k.Has(&c); has != tt.wantOK {
				t.Errorf("Has() = %v, want %v", has, tt.wantOK)
			}
		})
	}
}

func TestGetValueType(t *testing.T) {
	type config struct{ addr string }

	var c gocontainer.Container
	k := gocontainer.NewKey[config]("config")
	want := config{addr: ":8080"}
	k.Register(&c, want)

	if got := k.MustGet(&c); got != want {
		t.Errorf("MustGet() = %+v, want %+v", got, want)
	}
}

func TestRegisterInterfaceKey(t *testing.T) {
	tests := []struct {
		name     string
		register func(c *gocontainer.Container, k gocontainer.Key[Greeter])
		want     Greeter
		wantOK   bool
	}{
		{
			// A concrete value goes in under an interface key with no conversion at the call.
			name:     "implementation",
			register: func(c *gocontainer.Container, k gocontainer.Key[Greeter]) { k.Register(c, politeGreeter{"Hi"}) },
			want:     politeGreeter{"Hi"},
			wantOK:   true,
		},
		{
			name:     "nil interface",
			register: func(c *gocontainer.Container, k gocontainer.Key[Greeter]) { k.Register(c, nil) },
			wantOK:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c gocontainer.Container
			k := gocontainer.NewKey[Greeter]("greeter")
			tt.register(&c, k)

			got, ok := k.Get(&c)
			if got != tt.want || ok != tt.wantOK {
				t.Errorf("Get() = %v, %v, want %v, %v", got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestMustGet(t *testing.T) {
	var c gocontainer.Container
	db := &service{name: "db"}
	k := gocontainer.NewKey[*service]("db")
	k.Register(&c, db)

	if got := k.MustGet(&c); got != db {
		t.Errorf("MustGet() = %v, want %v", got, db)
	}
}

func TestMustGetPanics(t *testing.T) {
	tests := []struct {
		name string
		get  func(c *gocontainer.Container)
		want string
	}{
		{
			name: "not registered",
			get:  func(c *gocontainer.Container) { gocontainer.NewKey[*service]("db").MustGet(c) },
			want: `gocontainer: no *gocontainer_test.service registered as "db"`,
		},
		{
			name: "interface key",
			get:  func(c *gocontainer.Container) { gocontainer.NewKey[Greeter]("greeter").MustGet(c) },
			want: `gocontainer: no gocontainer_test.Greeter registered as "greeter"`,
		},
		{
			name: "deregistered",
			get: func(c *gocontainer.Container) {
				k := gocontainer.NewKey[int]("port")
				k.Register(c, 8080)
				k.Deregister(c)
				k.MustGet(c)
			},
			want: `gocontainer: no int registered as "port"`,
		},
		{
			name: "zero key",
			get:  func(c *gocontainer.Container) { gocontainer.Key[error]{}.MustGet(c) },
			want: `gocontainer: no error registered as ""`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c gocontainer.Container
			if got := wantPanic(t, func() { tt.get(&c) }); got != tt.want {
				t.Errorf("panic = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMustGetPanicsWithAnError(t *testing.T) {
	defer func() {
		if _, ok := recover().(error); !ok {
			t.Error("MustGet did not panic with an error")
		}
	}()
	var c gocontainer.Container
	gocontainer.NewKey[int]("port").MustGet(&c)
}

func TestZeroKey(t *testing.T) {
	var c gocontainer.Container
	var zero gocontainer.Key[*service]

	got := wantPanic(t, func() { zero.Register(&c, &service{}) })
	if want := "gocontainer: Register with a zero Key; make keys with NewKey"; got != want {
		t.Errorf("Register panic = %q, want %q", got, want)
	}
	if v, ok := zero.Get(&c); v != nil || ok {
		t.Errorf("Get() = %v, %v, want <nil>, false", v, ok)
	}
	if zero.Has(&c) {
		t.Error("Has() = true, want false")
	}
	// Every zero key is the same key, whatever its type: had Register gone through,
	// a zero Key[string] would find the *service.
	if v, ok := (gocontainer.Key[string]{}).Get(&c); v != "" || ok {
		t.Errorf("zero Key[string] Get() = %q, %v, want \"\", false", v, ok)
	}
	zero.Deregister(&c)
}

func TestDeregister(t *testing.T) {
	var c gocontainer.Container
	db, cache := gocontainer.NewKey[*service]("db"), gocontainer.NewKey[*service]("cache")
	cacheService := &service{name: "cache"}
	db.Register(&c, &service{name: "db"})
	cache.Register(&c, cacheService)

	db.Deregister(&c)
	db.Deregister(&c) // already gone: a no-op

	if db.Has(&c) {
		t.Error("db still registered after Deregister")
	}
	if got, ok := cache.Get(&c); got != cacheService || !ok {
		t.Errorf("cache Get() = %v, %v after deregistering db, want %v, true", got, ok, cacheService)
	}
}

func TestZeroContainer(t *testing.T) {
	var c gocontainer.Container
	k := gocontainer.NewKey[*service]("db")

	k.Deregister(&c)
	if v, ok := k.Get(&c); v != nil || ok {
		t.Errorf("Get() = %v, %v, want <nil>, false", v, ok)
	}
	if k.Has(&c) {
		t.Error("Has() = true, want false")
	}
}

// TestConcurrentRegister registers from many goroutines into a zero Container at once, while others read.
// v1 checked for a missing map outside its lock; here no write may be lost, and -race must stay quiet.
// TestConcurrentDeregister registers from many goroutines while another key is registered and deregistered:
// Deregister copies the registry too, and must not drop what was registered meanwhile.
func TestConcurrentDeregister(t *testing.T) {
	const n = 64
	var c gocontainer.Container
	keys := make([]gocontainer.Key[int], n)
	for i := range keys {
		keys[i] = gocontainer.NewKey[int](fmt.Sprint(i))
	}
	other := gocontainer.NewKey[int]("other")

	var wg sync.WaitGroup
	for i, k := range keys {
		wg.Go(func() { k.Register(&c, i) })
		wg.Go(func() {
			for range 16 {
				other.Register(&c, i)
				other.Deregister(&c)
			}
		})
	}
	wg.Wait()

	for i, k := range keys {
		if got, ok := k.Get(&c); got != i || !ok {
			t.Errorf("key %d: Get() = %v, %v, want %v, true", i, got, ok, i)
		}
	}
}

func TestConcurrentRegister(t *testing.T) {
	const n = 64
	var c gocontainer.Container
	keys := make([]gocontainer.Key[int], n)
	for i := range keys {
		keys[i] = gocontainer.NewKey[int](fmt.Sprint(i))
	}

	var wg sync.WaitGroup
	for i, k := range keys {
		wg.Go(func() { k.Register(&c, i) })
		wg.Go(func() { keys[n-1-i].Get(&c) })
	}
	wg.Wait()

	for i, k := range keys {
		if got, ok := k.Get(&c); got != i || !ok {
			t.Errorf("key %d: Get() = %v, %v, want %v, true", i, got, ok, i)
		}
	}
}

// TestConcurrentReplace reads a service while other goroutines replace and deregister its neighbours:
// a reader sees one of the registered values, never a missing or torn one.
func TestConcurrentReplace(t *testing.T) {
	var c gocontainer.Container
	k, other := gocontainer.NewKey[error]("err"), gocontainer.NewKey[int]("other")
	errA, errB := errors.New("a"), errors.New("b")
	k.Register(&c, errA)

	var wg sync.WaitGroup
	for i := range 32 {
		wg.Go(func() {
			if i%2 == 0 {
				k.Register(&c, errB)
			} else {
				k.Register(&c, errA)
			}
			other.Register(&c, i)
			other.Deregister(&c)
		})
		wg.Go(func() {
			if got := k.MustGet(&c); got != errA && got != errB {
				t.Errorf("MustGet() = %v, want a or b", got)
			}
		})
	}
	wg.Wait()
}
