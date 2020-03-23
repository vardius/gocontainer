🪣 gocontainer
================
[![Build Status](https://travis-ci.org/vardius/gocontainer.svg?branch=master)](https://travis-ci.org/vardius/gocontainer)
[![Go Report Card](https://goreportcard.com/badge/github.com/vardius/gocontainer)](https://goreportcard.com/report/github.com/vardius/gocontainer)
[![codecov](https://codecov.io/gh/vardius/gocontainer/branch/master/graph/badge.svg)](https://codecov.io/gh/vardius/gocontainer)
[![](https://godoc.org/github.com/vardius/gocontainer?status.svg)](https://pkg.go.dev/github.com/vardius/gocontainer)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/vardius/gocontainer/blob/master/LICENSE.md)

<img align="right" height="180px" src="https://github.com/vardius/gorouter/blob/master/website/src/static/img/logo.png?raw=true" alt="logo" />

gocontainer - Dependency Injection Container

📖 ABOUT
==================================================
Contributors:

* [Rafał Lorenz](http://rafallorenz.com)

Want to contribute ? Feel free to send pull requests!

Have problems, bugs, feature ideas?
We are using the github [issue tracker](https://github.com/vardius/gocontainer/issues) to manage them.

## 📚 Documentation

For __examples__ **visit [godoc#pkg-examples](http://godoc.org/github.com/vardius/gocontainer#pkg-examples)**

For **GoDoc** reference, **visit [pkg.go.dev](https://pkg.go.dev/github.com/vardius/gocontainer)**

[MustInvokeMany](https://godoc.org/github.com/vardius/gocontainer#example-package--MustInvokeMany)

🚏 HOW TO USE
==================================================

First file `main.go` simply gets the repository from the container and prints it
we use **MustInvoke** method to simply present the way where we keep type safety
```go
package main

import (
    "github.com/vardius/gocontainer/example/repository"
    "github.com/vardius/gocontainer"
)

func main() {
    gocontainer.MustInvoke("repository.mysql", func(r Repository) {
        fmt.Println(r)
    })
}
```
Our database implementation uses `init()` function to register db service
```go
package database

import (
    "fmt"
    "database/sql"

    "github.com/vardius/gocontainer"
)

func NewDatabase() *sql.DB {
    db, _ := sql.Open("mysql", "dsn")

    return db
}

func init() {
    db := gocontainer.MustGet("db")

    gocontainer.Register("db", NewDatabase())
}
```
Our repository accesses earlier on registered db service
and following the same patter uses `init()` function to register repository service within container
```go
package repository

import (
    "fmt"
    "database/sql"

    "github.com/vardius/gocontainer"
    _ "github.com/vardius/gocontainer/example/database"
)

type Repository interface {}

func NewRepository(db *sql.DB) Repository {
    return &mysqlRepository{db}
}

type mysqlRepository struct {
    db *sql.DB
}

func init() {
    db := gocontainer.MustGet("db")

    gocontainer.Register("repository.mysql", NewRepository(db.(*sql.DB)))
}
```
You can disable global container instance by setting `gocontainer.GlobalContainer` to `nil`.
This package allows you to create many containers.
```go
package main

import (
    "github.com/vardius/gocontainer/example/repository"
    "github.com/vardius/gocontainer"
)

func main() {
    // disable global container instance
    gocontainer.GlobalContainer = nil

    mycontainer := gocontainer.New()
    mycontainer.Register("test", 1)
}
```
Please check [GoDoc](http://godoc.org/github.com/vardius/gocontainer) for more methods and examples.

📜 [License](LICENSE.md)
-------

This package is released under the MIT license. See the complete license in the package
