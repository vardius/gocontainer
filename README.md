# 🪣 gocontainer

[![Test](https://github.com/vardius/gocontainer/actions/workflows/test.yml/badge.svg)](https://github.com/vardius/gocontainer/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/vardius/gocontainer/v2.svg)](https://pkg.go.dev/github.com/vardius/gocontainer/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/vardius/gocontainer/v2)](https://goreportcard.com/report/github.com/vardius/gocontainer/v2)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE.md)

A dependency injection container for Go, with typed keys.

- **Typed keys.** A `Key[T]` names a service and fixes its type: `Register` takes a `T`, and `Get` returns one.
  No `interface{}` or type assertions in your code, and no reflection to find a service.
- **Checked when you build.** Registering a value of the wrong type under a key doesn't compile.
- **Reads don't wait.** `Get` takes no lock and allocates nothing, so goroutines reading at once don't slow each other down.
- **No globals.** The zero `Container` is ready to use, and nothing runs in `init()`.

📖 **[Documentation](https://vardius.github.io/gocontainer/)** · [API reference](https://pkg.go.dev/github.com/vardius/gocontainer/v2) · [Examples](examples) · [Migrating from v1](MIGRATION.md)

## Install

```sh
go get github.com/vardius/gocontainer/v2
```

Requires Go 1.26 or newer.

## Keys

A key names a service and fixes its type. Declare each key once, usually next to the service it names, and share it:
every `NewKey` call makes a different key, even for the same name and type. The name only labels the key in panic messages.

```go
var (
	DB    = gocontainer.NewKey[*sql.DB]("db")
	Users = gocontainer.NewKey[UserStore]("users")
)
```

## Register and get

The zero `Container` is empty and ready to use. `Register` takes the key's type, or anything assignable to it,
so a concrete type goes in under an interface key with no conversion. `MustGet` returns the key's type, and panics if
nothing is registered; `Get` reports it instead.

```go
var c gocontainer.Container
DB.Register(&c, db)
Users.Register(&c, NewSQLUsers(DB.MustGet(&c))) // a *SQLUsers, registered as a UserStore

users := Users.MustGet(&c) // a UserStore; panics if none is registered
users, ok := Users.Get(&c) // or reports it: ok is false if none is
```

`Register` replaces whatever was registered under the key, `Deregister` removes it, and `Has` reports whether there is one.

## Testing with fakes

A key works with any number of containers. To test code that takes its services from a container,
register a fake under the same key in a container of the test's own:

```go
func TestRoutesWhenTheDatabaseIsDown(t *testing.T) {
	var c gocontainer.Container
	Users.Register(&c, downUsers{}) // a fake whose every call fails

	w := httptest.NewRecorder()
	routes(&c).ServeHTTP(w, httptest.NewRequest("GET", "/users/1", nil))
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", w.Code)
	}
}
```

See [examples/server](examples/server) for the whole program and its test.

## Benchmarks

Time and allocations per call, on an Apple M4 Max with Go 1.27.0, median of 5 runs. Both versions ran on the same harness,
with a container of 32 services. With 16 goroutines at once, the time is wall time per call across all of them,
so it measures throughput, not how long one call takes. Run v2's own suite with `go test -bench .`.

| Scenario | v1 | v2 |
| --- | ---: | ---: |
| Get | 5.8 ns · 0 allocs | **3.6 ns · 0 allocs** |
| Get, not registered | 4.8 ns · 0 allocs | **3.3 ns · 0 allocs** |
| MustGet | 6.4 ns · 0 allocs | **3.5 ns · 0 allocs** |
| Has | 6.2 ns · 0 allocs | **3.4 ns · 0 allocs** |
| Get from 16 goroutines at once | 128 ns · 0 allocs | **0.27 ns · 0 allocs** |
| v1 `Invoke`, v2 `Get` | 80.9 ns · 0 allocs | **3.6 ns · 0 allocs** |
| v1 `MustInvoke`, v2 `MustGet` | 71.0 ns · 0 allocs | **3.5 ns · 0 allocs** |
| Three services: v1 `MustInvokeMany`, v2 three `MustGet`s | 163 ns · 5 allocs | **10.2 ns · 0 allocs** |
| Register, replacing a service | **9.5 ns · 0 allocs** | 225 ns · 5 allocs |

v1 read under a read-write lock, whose shared counter every reader writes to, so readers on different cores queue for it.
v2 keeps the services in a map it never changes once published: a reader loads it atomically and looks the key up.
To change it, `Register` and `Deregister` copy the map and swap the copy in, so they cost more than v1's
and grow with the number of services. They suit services registered at startup, not a registry that changes on every request.

## Upgrading from v1

String ids become typed keys, and `Invoke` gives way to `Get`, which returns the type you want:

```go
// v1
gocontainer.Register("db", openDB())
gocontainer.MustInvoke("db", func(db *sql.DB) {
	// use db
})

// v2
var DB = gocontainer.NewKey[*sql.DB]("db")

var c gocontainer.Container
DB.Register(&c, openDB())
db := DB.MustGet(&c) // use db
```

[MIGRATION.md](MIGRATION.md) maps every v1 call to its v2 form and lists the behaviour changes.

## Contributing

Bug reports and pull requests are welcome: see [CONTRIBUTING.md](.github/CONTRIBUTING.md).

## License

MIT, see [LICENSE.md](LICENSE.md).
