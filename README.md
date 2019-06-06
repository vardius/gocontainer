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

```go
package main

import "github.com/vardius/gocontainer"

func main() {
    gocontainer.Register("test.service")

    service, ok := gocontainer.Get("test.service")
    if !ok {
        log.Fatal("Service not found")    
    }

    service.DoSomething()
}
```

License
-------

This package is released under the MIT license. See the complete license in the package:

[LICENSE](LICENSE.md)
