Vardius - gocontainer
================
[![Build Status](https://travis-ci.org/vardius/gocontainer.svg?branch=master)](https://travis-ci.org/vardius/gocontainer)
[![Go Report Card](https://goreportcard.com/badge/github.com/vardius/gocontainer)](https://goreportcard.com/report/github.com/vardius/gocontainer)
[![codecov](https://codecov.io/gh/vardius/gocontainer/branch/master/graph/badge.svg)](https://codecov.io/gh/vardius/gocontainer)
[![](https://godoc.org/github.com/vardius/gocontainer?status.svg)](http://godoc.org/github.com/vardius/gocontainer)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/vardius/gocontainer/blob/master/LICENSE.md)

gocontainer - Dependency Injection Container

ABOUT
==================================================
Contributors:

* [Rafał Lorenz](http://rafallorenz.com)

Want to contribute ? Feel free to send pull requests!

Have problems, bugs, feature ideas?
We are using the github [issue tracker](https://github.com/vardius/gocontainer/issues) to manage them.

HOW TO USE
==================================================

1. [GoDoc](http://godoc.org/github.com/vardius/gocontainer)
2. [Examples](http://godoc.org/github.com/vardius/gocontainer#pkg-examples)

3. [MustInvokeMany](https://godoc.org/github.com/vardius/gocontainer#MustInvokeMany)

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
You can disable global container instance by setting `gocontainer.InitializeGlobalContainer` to `false`.
This package allows you to create many containers.
```go
package main

import (
    "github.com/vardius/gocontainer/example/repository"
    "github.com/vardius/gocontainer"
)

// disable global container instance
// remember to do it outside function body
gocontainer.InitializeGlobalContainer = false

func main() {
    mycontainer := gocontainer.New()
    mycontainer.Register("test", 1)
}
```
Please check [GoDoc](http://godoc.org/github.com/vardius/gocontainer) for more methods and examples.

License
-------

This package is released under the MIT license. See the complete license in the package:

[LICENSE](LICENSE.md)
