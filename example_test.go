package gocontainer_test

import (
	"fmt"

	gocontainer "github.com/vardius/gocontainer"
)

func Example_new() {
	// disable global container instance
	// remember to do it outside function body
	gocontainer.InitializeGlobalContainer = false

	mycontainer := gocontainer.New()
	mycontainer.Register("test", 1)

	fmt.Println(mycontainer.MustGet("test"))
	// Output:
	// 1
}

func Example_register() {
	gocontainer.Register("test", 1)

	fmt.Println(gocontainer.MustGet("test"))
	// Output:
	// 1
}

func Example_deregister() {
	gocontainer.Register("test", 1)
	gocontainer.Deregister("test")

	fmt.Println(gocontainer.Get("test"))
	// Output:
	// <nil> false
}

func Example_has() {
	gocontainer.Register("test", 1)

	fmt.Println(gocontainer.Has("test"))
	// Output:
	// true
}

func Example_get() {
	gocontainer.Register("test", 1)
	o, ok := gocontainer.Get("test")

	fmt.Println(o)
	fmt.Println(ok)
	// Output:
	// 1
	// true
}

func Example_mustGet() {
	gocontainer.Register("test", 1)
	o := gocontainer.MustGet("test")

	fmt.Println(o)
	// Output:
	// 1
}

func Example_invoke() {
	gocontainer.Register("test", 1)
	gocontainer.Invoke("test", func(i int, ok bool) {
		fmt.Println(i)
		fmt.Println(ok)
	})

	// Output:
	// 1
	// true
}

func Example_mustInvoke() {
	gocontainer.Register("test", 1)
	gocontainer.MustInvoke("test", func(i int) {
		fmt.Println(i)
	})

	// Output:
	// 1
}

func Example_mustInvokeMany() {
	gocontainer.Register("test1", 1)
	gocontainer.Register("test2", 2)

	gocontainer.MustInvokeMany("test1", "test2")(func(x int, y int) {
		fmt.Println(x, "|", y)
	})

	// Output:
	// 1|2
}
