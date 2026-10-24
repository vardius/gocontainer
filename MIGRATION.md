# Migrating from v1 to v2

v2 keeps what v1 did: register services, get them back, check for them and remove them. Services are now
found by typed keys instead of string ids, so they come back as their own types, and the reflection-based
`Invoke` calls give way to `Get` and `MustGet`. The global container is gone: you declare your own.

## 1. Update the import path

```sh
go get github.com/vardius/gocontainer/v2
```

```go
import "github.com/vardius/gocontainer/v2"
```

v2 needs Go 1.26 or newer.

## 2. Before and after

```go
// v1
func init() {
	gocontainer.Register("repository.mysql", NewRepository(db))
}

func main() {
	gocontainer.MustInvoke("repository.mysql", func(r Repository) {
		fmt.Println(r)
	})
}
```

```go
// v2
var RepositoryKey = gocontainer.NewKey[Repository]("repository.mysql")

func main() {
	var c gocontainer.Container
	RepositoryKey.Register(&c, NewRepository(db))

	r := RepositoryKey.MustGet(&c)
	fmt.Println(r)
}
```

Each string id becomes a key, declared once with the type of the service it names. Declare it next to the service
and share the variable: `NewKey` makes a new key on every call, so two calls with the same name don't find each
other's service.

## 3. Every call

| v1 | v2 |
| --- | --- |
| `type Container interface` | `type Container struct`, used as `*gocontainer.Container` |
| `New()` | `var c gocontainer.Container`, or `new(gocontainer.Container)`: the zero value is ready to use |
| `c.Register(id, object)` | `key.Register(&c, object)`, with `key := gocontainer.NewKey[T](id)` declared once |
| `c.Deregister(id)` | `key.Deregister(&c)` |
| `c.Has(id) bool` | `key.Has(&c) bool` |
| `c.Get(id) (interface{}, bool)` | `key.Get(&c) (T, bool)` |
| `c.MustGet(id) interface{}` | `key.MustGet(&c) T` |
| `c.Invoke(id, func(o T, ok bool) { ... })` | `o, ok := key.Get(&c)` |
| `c.MustInvoke(id, func(o T) { ... })` | `o := key.MustGet(&c)` |
| `c.MustInvokeMany(idA, idB)(func(a A, b B) { ... })` | `a, b := keyA.MustGet(&c), keyB.MustGet(&c)` |
| `GlobalContainer` | a container you declare, such as `var Services gocontainer.Container` in a package of your own |
| `gocontainer.Register(id, object)` and the other package functions | the same key methods, given that container: `key.Register(&Services, object)` |
| `GlobalContainer = nil`, to turn the global off | nothing to do: there is no global |

A value of a concrete type goes in under an interface key as it is: `Register` takes the key's type or anything
assignable to it. A key of type `any` gives you v1's untyped behaviour back, if you need it.

## 4. Behaviour changes

- **Keys are compared by identity, not by name.** In v1 two equal strings were the same id. In v2 every
  `NewKey` call makes a different key, whatever its name and type, so share the key variable. The name only
  appears in panic messages.
- **Get reports a missing service.** v1's `Invoke` panicked inside reflect, with "reflect: Call using zero Value
  argument", when nothing was registered under the id, even though its func took an `ok`. v2's `Get` returns the
  zero value and false.
- **A registered nil comes back.** v1's `Invoke` and `MustInvoke` panicked inside reflect on a service registered
  as an untyped nil, `Register(id, nil)`. v2 returns nil, and `Get` reports it as registered.
- **Type mistakes don't compile.** v1 found a func of the wrong type, or a type assertion to the wrong type, only
  when it ran, and panicked. In v2, `Register` takes only the key's type and `Get` returns it.
- **Services are fetched when you call `MustGet`.** v1's `MustInvokeMany` fetched the services when it was called,
  not when the func it returned ran: it panicked before you had passed the func, and the func got the services
  registered at that moment, even if they had been replaced since.
- **`MustGet`'s panic names the type and the key.** v1 panicked with the string `Object <id> nof found within a container`.
  v2 panics with an error that reads `gocontainer: no *sql.DB registered as "db"`.
- **`Register` with the zero `Key` panics.** A key must come from `NewKey`. `Get`, `Has` and `Deregister` find nothing
  under the zero key.
- **No global, and nothing runs in `init()`.** v1 made `GlobalContainer` in its `init()`. With it set to nil,
  `MustInvokeMany` crashed with a nil pointer dereference instead of the "Global container instance is not
  initialized" panic the other functions gave. In v2 you declare your containers, and a test can use one of its own.
- **Reads are faster and writes are slower.** `Get`, `MustGet` and `Has` take no lock. `Register` and `Deregister`
  copy the registry, so they take time in proportion to the services in it: register your services at startup.

## 5. Registering in `init()`

v1's examples registered services in `init()` functions and imported packages with a blank import, only so that
their `init()` ran first. In v2, register services in `main`, in the order they need each other, or keep `init()`
with a container in a package of its own. A package that uses another package's key imports that package, and Go
initialises an imported package first, so the blank imports go away.
