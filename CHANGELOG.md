# Changelog

## v2.0.0 (2026-10-24)

Typed keys instead of string ids, and no reflection to find a service. See [MIGRATION.md](MIGRATION.md) for the before and after of every call.

### Added

- `Key[T]` and `NewKey`: a key names a service and fixes its type. `Register`, `Get`, `MustGet`, `Has` and `Deregister`
  are its methods, so a service goes in as a `T` and comes out as one, with no type assertion, and a value of the
  wrong type doesn't compile.
- The zero `Container` is ready to use.
- Runnable examples in `examples/`, each with a test: wiring services at startup, and an HTTP server tested with a fake.

### Changed

- Module path is `github.com/vardius/gocontainer/v2`; Go 1.26 or newer is required.
- Keys are compared by identity: every `NewKey` call makes a different key, even for the same name and type.
- `Get`, `MustGet` and `Has` take no lock and allocate nothing. From 16 goroutines at once, `Get` went from 128 ns to 0.27 ns;
  from one, from 5.8 ns to 3.6 ns. Allocation counts are guarded by tests.
- `Register` and `Deregister` copy the registry, so reads never wait for them: they take time in proportion to the services
  registered. Replacing one of 32 services went from 9.5 ns to 225 ns.
- `MustGet` panics with an error naming the type and the key: `gocontainer: no *sql.DB registered as "db"`.
- `Register` panics with the zero `Key`.

### Removed

- The `Container` interface and `New`: `Container` is a struct.
- `Invoke`, `MustInvoke` and `MustInvokeMany`: `Get` and `MustGet` return the type you want.
- `GlobalContainer`, the package functions that used it, and the package's `init()`.

### Fixed

- `Invoke` panicked inside reflect when nothing was registered under the id, and `Invoke` and `MustInvoke` did when
  an untyped nil was. `Get` returns the zero value and false, or the nil that was registered and true.
- `MustInvokeMany` crashed with a nil pointer dereference when `GlobalContainer` was nil.
