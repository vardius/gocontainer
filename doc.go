/*
Package gocontainer is a simple dependency injection container

Take the following example:
First file `main.go` simply gets the repository from the container and prints it
we use **MustInvoke** method to simply present the way where we keep type safety
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
Our database implementation uses `init()` function to register db service
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
Our repository accesses earlier on registered db service
and following the same patter uses `init()` function to register repository service within container
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
You can disable global container instance by setting gocontainer.InitializeGlobalContainer to false.
This package allows you to create many containers.
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
*/
package gocontainer
